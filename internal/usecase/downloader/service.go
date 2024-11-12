package downloader

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"hafiedh.com/downloader/internal/config"
	"hafiedh.com/downloader/internal/domain/entities"
	"hafiedh.com/downloader/internal/domain/repositories"
	"hafiedh.com/downloader/internal/infrastructure/imagekit"
	"hafiedh.com/downloader/internal/infrastructure/sse"
	"hafiedh.com/downloader/internal/pkg/constants"
)

const (
	UploadProviderImageKit = "imagekit"
	UploadProviderGoogle   = "google"
)

type (
	FileTransfer struct {
		imageKit         imagekit.ImageKitWrapper
		fileTransferRepo repositories.FileTransfer
		semaphore        chan struct{}
		sseHub           *sse.SSEHub
		Concurrency      int
	}
)

func NewService(fileTransferRepo repositories.FileTransfer, maxConcurrent int, imgKit imagekit.ImageKitWrapper, sseHub *sse.SSEHub) FileTransfer {
	return FileTransfer{
		fileTransferRepo: fileTransferRepo,
		Concurrency:      maxConcurrent,
		semaphore:        make(chan struct{}, maxConcurrent),
		imageKit:         imgKit,
		sseHub:           sseHub,
	}
}

func (f *FileTransfer) UploadFile(ctx context.Context, files map[string][]*multipart.FileHeader, values map[string][]string) (resp constants.DefaultResponse, err error) {
	// Use a longer timeout for large uploads
	uploadCtx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	resp = constants.DefaultResponse{
		Message: constants.MESSAGE_FAILED,
		Status:  http.StatusInternalServerError,
		Data:    struct{}{},
	}

	// Add request ID for better tracking
	requestID := uuid.New().String()
	logger := slog.With(
		"requestID", requestID,
		"totalFiles", countTotalFiles(files),
	)
	logger.InfoContext(uploadCtx, "starting upload batch")

	provider := config.GetString("multipart.uploadProvider")
	if !isValidProvider(provider) {
		logger.ErrorContext(uploadCtx, "invalid provider", "provider", provider)
		return resp, fmt.Errorf("invalid upload provider: %s", provider)
	}

	totalFiles := countTotalFiles(files)
	if totalFiles == 0 {
		return resp, fmt.Errorf("no files to upload")
	}

	// Add file size validation
	if err := f.validateFiles(files); err != nil {
		logger.ErrorContext(uploadCtx, "file size validation failed", "error", err)
		return resp, fmt.Errorf("file size validation failed: %w", err)
	}

	var wg sync.WaitGroup
	results := make(chan UploadResult, totalFiles)
	jobs := make(chan FileUploadJob, totalFiles)
	workerErrors := make(chan error, f.Concurrency)

	// Monitor worker health
	healthCtx, healthCancel := context.WithCancel(uploadCtx)
	defer healthCancel()

	go f.monitorWorkerHealth(healthCtx, workerErrors, logger)

	// Start workers with improved error handling
	for i := 0; i < f.Concurrency; i++ {
		workerID := i
		go func() {
			if err := f.worker(uploadCtx, jobs, results, provider); err != nil {
				select {
				case workerErrors <- fmt.Errorf("worker %d error: %w", workerID, err):
				default:
					logger.ErrorContext(uploadCtx, "worker error",
						"workerID", workerID,
						"error", err)
				}
			}
		}()
	}

	// Dispatch jobs with backpressure
	go func() {
		defer close(jobs)
		for key, fileHeaders := range files {
			for _, fileHeader := range fileHeaders {
				select {
				case <-uploadCtx.Done():
					return
				default:
					wg.Add(1)

					// Add retry logic for job submission
					err := retry.Do(
						func() error {
							select {
							case jobs <- FileUploadJob{
								Key:        key,
								FileHeader: fileHeader,
								Result:     results,
								RequestID:  requestID,
							}:
								return nil
							case <-uploadCtx.Done():
								wg.Done()
								return uploadCtx.Err()
							case <-time.After(5 * time.Second):
								return fmt.Errorf("timeout queuing job")
							}
						},
						retry.Attempts(3),
						retry.Delay(time.Second),
						retry.Context(uploadCtx),
					)

					if err != nil {
						logger.ErrorContext(uploadCtx, "failed to queue job",
							"filename", fileHeader.Filename,
							"error", err)
						wg.Done()
					}
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var uploadErrors []error
	successCount := 0
	uploadStartTime := time.Now()

	// Process results with timeout
	for {
		select {
		case <-uploadCtx.Done():
			logger.ErrorContext(uploadCtx, "upload operation cancelled or timed out",
				"duration", time.Since(uploadStartTime),
				"successCount", successCount)
			uploadErrors = append(uploadErrors, fmt.Errorf("upload operation cancelled or timed out: %w", uploadCtx.Err()))
			goto HANDLE_RESULTS

		case err := <-workerErrors:
			logger.ErrorContext(uploadCtx, "critical worker error",
				"duration", time.Since(uploadStartTime),
				"successCount", successCount,
				"error", err)
			uploadErrors = append(uploadErrors, fmt.Errorf("worker error: %w", err))
			goto HANDLE_RESULTS

		case result, ok := <-results:
			if !ok {
				goto HANDLE_RESULTS
			}
			if result.Error != nil {
				if isContextCanceled(result.Error) {
					logger.WarnContext(uploadCtx, "upload cancelled for file",
						"filename", result.FileName,
						"error", result.Error)
					uploadErrors = append(uploadErrors, fmt.Errorf("upload cancelled for %s: %w", result.FileName, result.Error))
				} else {
					logger.ErrorContext(uploadCtx, "upload failed for file",
						"filename", result.FileName,
						"error", result.Error)
					uploadErrors = append(uploadErrors, fmt.Errorf("%s: %w", result.FileName, result.Error))
				}
			} else {
				successCount++
				logger.InfoContext(uploadCtx, "file upload successful",
					"filename", result.FileName,
					"url", result.URL)
			}
			wg.Done()
		}
	}

HANDLE_RESULTS:
	logger.InfoContext(uploadCtx, "upload batch completed",
		"duration", time.Since(uploadStartTime),
		"successCount", successCount,
		"errorCount", len(uploadErrors))

	if len(uploadErrors) > 0 {
		if successCount == 0 {
			resp.Message = constants.MESSAGE_FAILED
			resp.Status = http.StatusInternalServerError
			resp.Errors = formatErrors(uploadErrors)
			return resp, fmt.Errorf("all uploads failed: %v", uploadErrors[0])
		}
		resp.Message = constants.MESSAGE_PARTIAL_SUCCESS
		resp.Status = http.StatusPartialContent
		resp.Errors = formatErrors(uploadErrors)
		return resp, nil
	}
	resp.Message = constants.MESSAGE_SUCCESS
	resp.Status = http.StatusOK
	return resp, nil
}

func (f *FileTransfer) worker(ctx context.Context, jobs <-chan FileUploadJob, results chan<- UploadResult, provider string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job, ok := <-jobs:
			if !ok {
				return nil
			}
			select {
			case f.semaphore <- struct{}{}:
				result := f.processFile(ctx, job, provider)
				select {
				case results <- result:
					<-f.semaphore
				case <-ctx.Done():
					<-f.semaphore
					return ctx.Err()
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (f *FileTransfer) processFile(ctx context.Context, job FileUploadJob, provider string) UploadResult {
	file, err := job.FileHeader.Open()
	if err != nil {
		return UploadResult{
			FileName: job.FileHeader.Filename,
			Error:    fmt.Errorf("failed to open file: %w", err),
		}
	}
	defer file.Close()

	bufReader := bufio.NewReaderSize(file, 1024*1024)

	pr := NewProgressReader(bufReader, job.FileHeader.Size, func(progress UploadProgress) {
		if ctx.Err() != nil {
			return
		}
		progressMsg := constants.DefaultResponse{
			Message: constants.MESSAGE_PROGRESSION,
			Status:  http.StatusPartialContent,
			Data:    progress,
		}
		f.sseHub.BroadcastMessage(progressMsg)
	})

	var metaData *entities.MetaData

	switch provider {
	case UploadProviderImageKit:
		var wrapperResp imagekit.UploadResponse
		err := retry.Do(
			func() error {
				var uploadErr error
				wrapperResp, uploadErr = f.imageKit.UploadFile(ctx, pr, job.Key, job.FileHeader.Filename)
				return uploadErr
			},
			retry.Attempts(3),
			retry.Delay(time.Second),
			retry.DelayType(retry.BackOffDelay),
			retry.Context(ctx),
			retry.OnRetry(func(n uint, err error) {
				slog.WarnContext(ctx, "retrying upload",
					"filename", job.FileHeader.Filename,
					"attempt", n+1,
					"error", err)
			}),
		)

		if err != nil {
			if isContextCanceled(err) {
				return UploadResult{
					FileName: job.FileHeader.Filename,
					Error:    fmt.Errorf("upload cancelled: %w", err),
				}
			}
			return UploadResult{
				FileName: job.FileHeader.Filename,
				Error:    fmt.Errorf("imagekit upload failed after retries: %w", err),
			}
		}

		metaData = &entities.MetaData{
			URL:       wrapperResp.URL,
			FileName:  job.FileHeader.Filename,
			FileSize:  wrapperResp.Size,
			FileID:    wrapperResp.FileID,
			Extension: strings.ToLower(strings.TrimPrefix(filepath.Ext(job.FileHeader.Filename), ".")),
			Status:    "UPLOADED",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

	case UploadProviderGoogle:
		return UploadResult{
			FileName: job.FileHeader.Filename,
			Error:    fmt.Errorf("google upload provider not implemented"),
		}

	default:
		return UploadResult{
			FileName: job.FileHeader.Filename,
			Error:    fmt.Errorf("unsupported upload provider: %s", provider),
		}
	}

	err = retry.Do(
		func() error {
			return f.fileTransferRepo.SaveUpload(ctx, metaData)
		},
		retry.Attempts(3),
		retry.Delay(time.Second),
		retry.Context(ctx),
	)

	if err != nil {
		return UploadResult{
			FileName: job.FileHeader.Filename,
			Error:    fmt.Errorf("failed to save metadata"),
		}
	}

	return UploadResult{
		FileName: job.FileHeader.Filename,
		URL:      metaData.URL,
		Error:    nil,
	}
}

func countTotalFiles(files map[string][]*multipart.FileHeader) int {
	total := 0
	for _, headers := range files {
		total += len(headers)
	}
	return total
}

func formatErrors(errors []error) []string {
	formatted := make([]string, len(errors))
	for i, err := range errors {
		formatted[i] = err.Error()
	}
	return formatted
}

func isValidProvider(provider string) bool {
	switch provider {
	case UploadProviderImageKit, UploadProviderGoogle:
		return true
	default:
		return false
	}
}

func isContextCanceled(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func (f *FileTransfer) monitorWorkerHealth(ctx context.Context, workerErrors <-chan error, logger *slog.Logger) {
	errorCount := 0
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-workerErrors:
			errorCount++
			logger.ErrorContext(ctx, "worker error detected",
				"errorCount", errorCount,
				"error", err)

			// Consider implementing recovery logic here
			if errorCount > f.Concurrency/2 {
				logger.ErrorContext(ctx, "critical error threshold reached",
					"errorCount", errorCount,
					"concurrency", f.Concurrency)
			}
		}
	}
}

func (f *FileTransfer) validateFiles(files map[string][]*multipart.FileHeader) error {
	const maxFileSize = 100 * 1024 * 1024 // 100MB

	for _, fileHeaders := range files {
		for _, header := range fileHeaders {
			if header.Size > maxFileSize {
				return fmt.Errorf("file %s exceeds maximum size of 100MB", header.Filename)
			}
		}
	}
	return nil
}

func (f *FileTransfer) PresignedFile(ctx context.Context, id int64) (resp FilePresignedUrlResponse, err error) {
	resp = FilePresignedUrlResponse{
		DefaultResponse: constants.DefaultResponse{
			Message: constants.MESSAGE_FAILED,
			Status:  http.StatusInternalServerError,
			Errors:  make([]string, 0),
		},
	}

	metaData, err := f.fileTransferRepo.FindByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find metadata", "error", err)
		err = fmt.Errorf("failed to find metadata")
		return
	}

	var presignedURL string
	switch {
	case strings.HasPrefix(metaData.URL, config.GetString("imagekit.baseURL")):
		presignedURL, err = f.imageKit.PresignURL(ctx, metaData.URL)
		if err != nil {
			slog.ErrorContext(ctx, "failed to presign url", "error", err)
			err = fmt.Errorf("failed to presign url: %w", err)
			return
		}
	case strings.HasPrefix(metaData.URL, "https://storage.googleapis.com"):
		return resp, fmt.Errorf("google presign url not implemented")
	default:
		return resp, fmt.Errorf("unsupported url: %s", metaData.URL)
	}
	resp = FilePresignedUrlResponse{
		Data: FilePresigned{
			Url: presignedURL,
		},
		DefaultResponse: constants.DefaultResponse{
			Message: constants.MESSAGE_SUCCESS,
			Status:  http.StatusOK,
			Errors:  make([]string, 0),
		},
	}
	return

}
