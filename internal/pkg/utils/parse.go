package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	FileValidation struct {
		MaxFileSize   int64
		AllowedKeys   []string
		AllowedValues []string
		AllowedTypes  []string
		MaxNameLength int
		BlockedExtens []string
	}
	FileValidator struct {
		config *FileValidation
	}
)

func newDefaultFileValidation() *FileValidation {
	return &FileValidation{
		MaxFileSize:   5 * 1024 * 1024, // 5MB
		AllowedKeys:   []string{"file-1", "file-2", "file-3", "provider"},
		AllowedValues: []string{"imagekit", "google"},
		AllowedTypes:  []string{"image/jpeg", "image/png", "image/gif", "application/pdf", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/msword"},
		MaxNameLength: 255,
		BlockedExtens: []string{".exe", ".bat", ".cmd", ".sh", ".php", ".js"},
	}
}

func NewFileValidator(config *FileValidation) *FileValidator {
	if config == nil {
		config = newDefaultFileValidation()
	}
	return &FileValidator{config: config}
}

func (v *FileValidator) ValidateFileHeader(fileHeader *multipart.FileHeader) error {
	if len(fileHeader.Filename) > v.config.MaxNameLength {
		return fmt.Errorf("filename too long (max %d characters)", v.config.MaxNameLength)
	}
	if fileHeader.Size > v.config.MaxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d bytes", fileHeader.Size, v.config.MaxFileSize)
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	for _, blocked := range v.config.BlockedExtens {
		if ext == blocked {
			return fmt.Errorf("file extension %s is not allowed", ext)
		}
	}
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()
	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to read file header: %w", err)
	}
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = http.DetectContentType(buffer)
	}
	for _, allowedType := range v.config.AllowedTypes {
		if strings.HasPrefix(mimeType, allowedType) {
			return nil
		}
	}
	return fmt.Errorf("file type %s is not allowed", mimeType)
}

func (v *FileValidator) ValidateKey(key string) error {
	for _, allowedKey := range v.config.AllowedKeys {
		if key == allowedKey {
			return nil
		}
	}
	return fmt.Errorf("key %s is not allowed", key)
}

func (v *FileValidator) ValidateValue(key, value string) error {
	for _, allowedValue := range v.config.AllowedValues {
		if value == allowedValue {
			return nil
		}
	}
	return fmt.Errorf("value %s is not allowed for key %s", value, key)
}

func ParseMultipartForm(c echo.Context, validationConfig FileValidator) (files map[string][]*multipart.FileHeader, values map[string][]string, err error) {
	if err = c.Request().ParseMultipartForm(20 << 20); err != nil {
		return nil, nil, fmt.Errorf("Error:Field validation for parse multipart form")
	}
	files, values = c.Request().MultipartForm.File, c.Request().MultipartForm.Value
	if len(files) < 1 && len(values) < 1 {
		return nil, nil, fmt.Errorf("Error:Field validation for form file. Empty form")
	}

	for key, value := range values {
		if err = validationConfig.ValidateKey(key); err != nil {
			return
		}
		for _, val := range value {
			if err = validationConfig.ValidateValue(key, val); err != nil {
				return
			}
		}
		break
	}

	for key, valuess := range files {
		if err = validationConfig.ValidateKey(key); err != nil {
			return
		}
		for _, fileHeader := range valuess {
			if err = validationConfig.ValidateFileHeader(fileHeader); err != nil {
				return
			}
		}
	}
	return
}

func ParseInt16FromString(s string) (i int16, err error) {
	int64val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	i = int16(int64val)
	return
}
