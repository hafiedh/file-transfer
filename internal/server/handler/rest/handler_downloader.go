package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"hafiedh.com/downloader/internal/config"
	"hafiedh.com/downloader/internal/pkg/constants"
	"hafiedh.com/downloader/internal/pkg/utils"
	"hafiedh.com/downloader/internal/usecase/downloader"
)

type (
	downloadHandler struct {
		downloaderService downloader.FileTransfer
	}
)

func NewDownloaderHandler(downloaderService downloader.FileTransfer) *downloadHandler {
	if downloaderService.Concurrency == 0 {
		panic("downloaderService is not implemented yet")
	}
	return &downloadHandler{
		downloaderService: downloaderService,
	}
}

func (h *downloadHandler) Uploader(ec echo.Context) (err error) {

	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ec.Request().Context(), "panic occurred", "error", r)
			ec.JSON(http.StatusInternalServerError, constants.DefaultResponse{
				Message: "internal server error",
				Status:  http.StatusInternalServerError,
				Data:    struct{}{},
				Errors:  []string{"internal server error"},
			})
		}
	}()
	ctx := ec.Request().Context()

	config := utils.FileValidation{
		MaxFileSize:   config.GetInt64("multipart.maxFileSize"),
		AllowedKeys:   strings.Split(config.GetString("multipart.allowed_keys"), ","),
		AllowedValues: strings.Split(config.GetString("multipart.allowed_values"), ","),
		AllowedTypes:  strings.Split(config.GetString("multipart.allowed_types"), ","),
		MaxNameLength: config.GetInt("multipart.maxNameLength"),
		BlockedExtens: strings.Split(config.GetString("multipart.blockedExtens"), ","),
	}
	var validationConfig *utils.FileValidator
	if config.AllowedKeys == nil && config.AllowedValues == nil {
		validationConfig = utils.NewFileValidator(nil)
	} else {
		validationConfig = utils.NewFileValidator(&config)
	}

	files, values, err := utils.ParseMultipartForm(ec, *validationConfig)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse multipart form", "error", err)
		return ec.JSON(http.StatusBadRequest, constants.DefaultResponse{
			Message: "failed to parse multipart form",
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Errors:  []string{err.Error()},
		})
	}

	resp, err := h.downloaderService.UploadFile(ctx, files, values)
	if err != nil {
		slog.ErrorContext(ctx, "failed to upload file", "error", err)
		return ec.JSON(http.StatusInternalServerError, constants.DefaultResponse{
			Message: "failed to upload file",
			Status:  http.StatusInternalServerError,
			Data:    struct{}{},
			Errors:  []string{err.Error()},
		})
	}

	return ec.JSON(http.StatusOK, constants.DefaultResponse{
		Message: "success",
		Status:  http.StatusOK,
		Data:    resp.Data,
		Errors:  make([]string, 0),
	})
}

func (h *downloadHandler) Presigned(ec echo.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ec.Request().Context(), "panic occurred", "error", r)
			ec.JSON(http.StatusInternalServerError, constants.DefaultResponse{
				Message: "internal server error",
				Status:  http.StatusInternalServerError,
				Data:    struct{}{},
				Errors:  []string{"internal server error"},
			})
		}
	}()
	ctx := ec.Request().Context()

	id := ec.Param("id")
	if id == "" {
		slog.ErrorContext(ctx, "id is required")
		return ec.JSON(http.StatusBadRequest, constants.DefaultResponse{
			Message: "id is required",
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Errors:  []string{"id is required"},
		})
	}

	fileID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse id", "error", err)
		err = fmt.Errorf("failed to parse id: %w", err)
		return ec.JSON(http.StatusBadRequest, constants.DefaultResponse{
			Message: "failed to parse id",
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Errors:  []string{err.Error()},
		})
	}

	resp, err := h.downloaderService.PresignedFile(ctx, fileID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to presign file", "error", err)
		return ec.JSON(http.StatusInternalServerError, constants.DefaultResponse{
			Message: "failed to presign file",
			Status:  http.StatusInternalServerError,
			Data:    struct{}{},
			Errors:  []string{err.Error()},
		})
	}
	return ec.JSON(http.StatusOK, resp)
}
