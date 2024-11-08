package imagekit

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	ik "github.com/imagekit-developer/imagekit-go"
	ikUrl "github.com/imagekit-developer/imagekit-go/url"
	"hafiedh.com/downloader/internal/config"
	"hafiedh.com/downloader/internal/pkg/utils"
)

type (
	ImageKitWrapper interface {
		UploadFile(ctx context.Context, file io.Reader, key, fileName string) (UploadFileResp UploadResponse, err error)
		PresignURL(ctx context.Context, url string) (presignedURL string, err error)
	}

	imageKitWrapper struct {
		basicAuth  string
		privateKey string
		publicKey  string
	}
)

func NewImageKitWrapper() ImageKitWrapper {
	privateKey := config.GetString("imagekit.privateKey")
	publicKey := config.GetString("imagekit.publicKey")
	basicAuth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(privateKey)))
	return &imageKitWrapper{
		privateKey: privateKey,
		basicAuth:  basicAuth,
		publicKey:  publicKey,
	}
}

func (i *imageKitWrapper) UploadFile(ctx context.Context, file io.Reader, key, fileName string) (UploadFileResp UploadResponse, err error) {
	path := fmt.Sprintf("%s%s", config.GetString("imagekit.baseuploadURL"), config.GetString("imagekit.uploadPath"))

	if path == "" {
		slog.ErrorContext(ctx, "imagekit base URL is empty")
		err = fmt.Errorf("imagekit base URL is empty")
		return
	}
	uploadReq := UploadRequest{
		FileName:          fileName,
		UseUniqueFileName: "true",
		IsPrivateFile:     "true",
		Folder:            config.GetString("imagekit.uploadFolder"),
	}
	token, err := i.generateToken(uploadReq)

	if err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to generate token", "error", err)
		err = fmt.Errorf("failed to generate token: %w", err)
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to create form file", "error", err)
		err = fmt.Errorf("failed to create form file: %w", err)
		return
	}

	if _, err = io.Copy(part, file); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to copy file to part", "error", err)
		err = fmt.Errorf("failed to copy file to part: %w", err)
		return
	}

	if err = writer.WriteField("fileName", fileName); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to write field", "error", err)
		err = fmt.Errorf("failed to write field: %w", err)
		return
	}

	if err = writer.WriteField("token", token); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to write field", "error", err)
		err = fmt.Errorf("failed to write field: %w", err)
		return
	}

	if err = writer.WriteField("useUniqueFileName", "true"); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to write field", "error", err)
		err = fmt.Errorf("failed to write field: %w", err)
		return
	}

	if err = writer.WriteField("folder", uploadReq.Folder); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to write field", "error", err)
		err = fmt.Errorf("failed to write field: %w", err)
		return
	}

	if err = writer.WriteField("isPrivateFile", uploadReq.IsPrivateFile); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to write field", "error", err)
		err = fmt.Errorf("failed to write field: %w", err)
		return
	}

	if err = writer.Close(); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to close writer", "error", err)
		err = fmt.Errorf("failed to close writer: %w", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, body)
	if err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to create new HTTP request", "error", err)
		err = fmt.Errorf("failed to create new HTTP request: %w", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", i.basicAuth)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	uploadResp, err := client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to perform HTTP request", "error", err)
		err = fmt.Errorf("failed to perform HTTP request: %w", err)
		return
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK && uploadResp.StatusCode != http.StatusCreated {
		var badResponse BadResponse
		slog.ErrorContext(ctx, "[ImageKit] failed to upload file", "status", uploadResp.Status)
		if err = json.NewDecoder(uploadResp.Body).Decode(&badResponse); err != nil {
			slog.ErrorContext(ctx, "[ImageKit] failed to decode response body", "error", err)
			err = fmt.Errorf("failed to decode response body: %w", err)
			return
		}
		err = fmt.Errorf("failed to upload file: %s", badResponse.Message)
		return
	}

	var uploadResponse UploadResponse
	if err = json.NewDecoder(uploadResp.Body).Decode(&uploadResponse); err != nil {
		slog.ErrorContext(ctx, "[ImageKit] failed to decode upload response", "error", err)
		err = fmt.Errorf("failed to decode upload response: %w", err)
		return
	}

	return uploadResponse, nil
}

func (i *imageKitWrapper) generateToken(req UploadRequest) (token string, err error) {
	iat := time.Now().Unix()
	exp := time.Now().Add(10 * time.Minute).Unix()

	request := utils.UploadRequest{
		FileName:          req.FileName,
		UseUniqueFileName: req.UseUniqueFileName,
		Folder:            req.Folder,
		IsPrivateFile:     req.IsPrivateFile,
		Iat:               iat,
		Exp:               exp,
	}
	token, err = utils.ImageKitJwtSign(request)
	if err != nil {
		return
	}
	return
}

func (i *imageKitWrapper) PresignURL(ctx context.Context, url string) (presignedURL string, err error) {
	newIk := ik.NewFromParams(
		ik.NewParams{
			PrivateKey:  i.privateKey,
			PublicKey:   i.publicKey,
			UrlEndpoint: config.GetString("imagekit.urlEndpoint"),
		},
	)
	presignedURL, err = newIk.Url(
		ikUrl.UrlParam{
			Path:          strings.ReplaceAll(url, fmt.Sprintf("%s/", config.GetString("imagekit.urlEndpoint")), ""),
			Signed:        true,
			ExpireSeconds: 3600,
			UnixTime: func() int64 {
				return time.Now().Unix()
			},
		},
	)

	return
}
