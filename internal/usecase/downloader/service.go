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
	uploadCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resp = constants.DefaultResponse{
		Message: constants.MESSAGE_FAILED,
		Status:  http.StatusInternalServerError,
		Data:    struct{}{},
	}

	provider := config.GetString("multipart.uploadProvider")
	if !isValidProvider(provider) {
		return resp, fmt.Errorf("invalid upload provider: %s", provider)
	}

	totalFiles := countTotalFiles(files)
	if totalFiles == 0 {
		return resp, fmt.Errorf("no files to upload")
	}

	var wg sync.WaitGroup
	results := make(chan UploadResult, totalFiles)
	jobs := make(chan FileUploadJob, totalFiles)

	workerErrors := make(chan error, f.Concurrency)

	for i := 0; i < f.Concurrency; i++ {
		go func() {
			if err := f.worker(uploadCtx, jobs, results, provider); err != nil {
				select {
				case workerErrors <- err:
				default:
					// Channel is full, log the error
					slog.ErrorContext(uploadCtx, "worker error", "error", err)
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for key, fileHeaders := range files {
			for _, fileHeader := range fileHeaders {
				select {
				case <-uploadCtx.Done():
					return
				default:
					wg.Add(1)
					select {
					case jobs <- FileUploadJob{
						Key:        key,
						FileHeader: fileHeader,
						Result:     results,
					}:
					case <-uploadCtx.Done():
						wg.Done()
						return
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

	for {
		select {
		case <-uploadCtx.Done():
			uploadErrors = append(uploadErrors, fmt.Errorf("upload operation cancelled or timed out: %w", uploadCtx.Err()))
			goto HANDLE_RESULTS

		case err := <-workerErrors:
			uploadErrors = append(uploadErrors, fmt.Errorf("worker error: %w", err))
			goto HANDLE_RESULTS

		case result, ok := <-results:
			if !ok {
				goto HANDLE_RESULTS
			}
			if result.Error != nil {
				if isContextCanceled(result.Error) {
					uploadErrors = append(uploadErrors, fmt.Errorf("upload cancelled for %s: %w", result.FileName, result.Error))
				} else {
					uploadErrors = append(uploadErrors, fmt.Errorf("%s: %w", result.FileName, result.Error))
				}
			} else {
				successCount++
			}
			wg.Done()
		}
	}

HANDLE_RESULTS:
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
			FileSize:  int(wrapperResp.Size),
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
			Error:    fmt.Errorf("failed to save metadata: %w", err),
		}
	}

	return UploadResult{
		FileName: job.FileHeader.Filename,
		URL:      metaData.URL,
		Error:    nil,
	}
}

func (f *FileTransfer) DownloadFile(ctx context.Context, url, destinationPath string) (resp constants.DefaultResponse, err error) {
	f.semaphore <- struct{}{}
	defer func() {
		<-f.semaphore
	}()
	return
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
