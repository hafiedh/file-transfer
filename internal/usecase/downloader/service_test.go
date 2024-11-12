package downloader

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hafiedh.com/downloader/internal/config"
	"hafiedh.com/downloader/internal/domain/entities"
	"hafiedh.com/downloader/internal/infrastructure/imagekit"
	"hafiedh.com/downloader/internal/infrastructure/sse"
	"hafiedh.com/downloader/internal/pkg/constants"
)

type MockFileTransferRepo struct {
	mock.Mock
}

func (m *MockFileTransferRepo) SaveUpload(ctx context.Context, metaData *entities.MetaData) error {
	args := m.Called(ctx, metaData)
	return args.Error(0)
}

func (m *MockFileTransferRepo) SaveDownload(ctx context.Context, metaData *entities.MetaData) error {
	args := m.Called(ctx, metaData)
	return args.Error(0)
}

func (m *MockFileTransferRepo) FindByID(ctx context.Context, id int64) (entities.MetaData, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.MetaData), args.Error(1)
}

type MockImageKitWrapper struct {
	mock.Mock
}

func (m *MockImageKitWrapper) UploadFile(ctx context.Context, file io.Reader, fileName string, fileType string) (imagekit.UploadResponse, error) {
	args := m.Called(ctx, file, fileName, fileType)
	return args.Get(0).(imagekit.UploadResponse), args.Error(1)
}

func (m *MockImageKitWrapper) PresignURL(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

type MockFile struct {
	mock.Mock
}

func (m *MockFile) Open() (multipart.File, error) {
	args := m.Called()
	return args.Get(0).(multipart.File), args.Error(1)
}

func TestUploadFile(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	mockSSEHub := new(sse.SSEHub)

	service := NewService(mockRepo, 2, mockImageKit, mockSSEHub)

	files := map[string][]*multipart.FileHeader{
		"file1": {
			{
				Filename: "test1.jpg",
				Size:     1024,
				Header:   textproto.MIMEHeader{},
			},
		},
	}

	values := map[string][]string{}

	config.Set("multipart.uploadProvider", UploadProviderImageKit)

	mockImageKit.On("UploadFile", mock.Anything, mock.Anything, "test1.jpg").Return(imagekit.UploadResponse{
		URL:    "http://example.com/test1.jpg",
		Size:   1024,
		FileID: "file1",
	}, nil)

	mockRepo.On("SaveUpload", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	resp, err := service.UploadFile(ctx, files, values)

	assert.NoError(t, err)
	assert.Equal(t, constants.MESSAGE_SUCCESS, resp.Message)
	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestUploadFile_InvalidProvider(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	mockSSEHub := new(sse.SSEHub)

	service := NewService(mockRepo, 2, mockImageKit, mockSSEHub)

	files := map[string][]*multipart.FileHeader{}

	values := map[string][]string{}

	config.Set("multipart.uploadProvider", "invalidProvider")

	ctx := context.Background()
	resp, err := service.UploadFile(ctx, files, values)

	assert.Error(t, err)
	assert.Equal(t, constants.MESSAGE_FAILED, resp.Message)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
}

func TestPresignedFile(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	mockSSEHub := new(sse.SSEHub)

	service := NewService(mockRepo, 2, mockImageKit, mockSSEHub)

	metaData := &entities.MetaData{
		URL: "http://example.com/test1.jpg",
	}

	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(*metaData, nil)
	mockImageKit.On("PresignURL", mock.Anything, "http://example.com/test1.jpg").Return("http://example.com/presigned/test1.jpg", nil)

	ctx := context.Background()
	resp, err := service.PresignedFile(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, constants.MESSAGE_SUCCESS, resp.Message)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, "http://example.com/presigned/test1.jpg", resp.Data.Url)
}

func TestPresignedFile_NotFound(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	mockSSEHub := new(sse.SSEHub)

	service := NewService(mockRepo, 2, mockImageKit, mockSSEHub)

	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(entities.MetaData{}, errors.New("failed to find metadata"))

	ctx := context.Background()
	resp, err := service.PresignedFile(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, constants.MESSAGE_FAILED, resp.Message)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
}

func TestFileTransfer_validateFiles(t *testing.T) {
	service := NewService(nil, 2, nil, nil)

	fileContent := strings.Repeat("a", 101*1024*1024) // 101MB
	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     int64(len(fileContent)),
	}

	files := map[string][]*multipart.FileHeader{
		"file": {fileHeader},
	}

	err := service.validateFiles(files)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum size of 100MB")
}

func TestFileTransfer_UploadFile_UploadFailed(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	sseHub := sse.NewSSEHub()

	service := NewService(mockRepo, 2, mockImageKit, sseHub)

	fileContent := "test file content"
	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     int64(len(fileContent)),
	}

	files := map[string][]*multipart.FileHeader{
		"file": {fileHeader},
	}

	mockImageKit.On("UploadFile", mock.Anything, mock.Anything, "test.txt", "image/jpeg").Return(imagekit.UploadResponse{}, errors.New("failed to upload file"))

	values := map[string][]string{}

	resp, err := service.UploadFile(context.Background(), files, values)
	assert.Error(t, err)
	assert.Equal(t, constants.MESSAGE_FAILED, resp.Message)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
}

func TestFileTransfer_UploadFile_ContextCancelled(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	sseHub := sse.NewSSEHub()

	service := NewService(mockRepo, 2, mockImageKit, sseHub)

	fileContent := "test file content"
	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     int64(len(fileContent)),
	}

	files := map[string][]*multipart.FileHeader{
		"file": {fileHeader},
	}

	config.Set("multipart.uploadProvider", UploadProviderImageKit)

	mockImageKit.On("UploadFile", mock.Anything, mock.Anything, "test.txt", "image/jpeg").Return(imagekit.UploadResponse{
		URL:    "http://example.com/test.txt",
		Size:   int64(len(fileContent)),
		FileID: "fileID",
	}, nil)

	mockRepo.On("SaveUpload", mock.Anything, mock.Anything).Return(nil)
	values := map[string][]string{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := service.UploadFile(ctx, files, values)
	assert.Error(t, err)
	assert.Equal(t, constants.MESSAGE_FAILED, resp.Message)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
}

func TestProcessFile_FailedToOpenFile(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	mockSSEHub := new(sse.SSEHub)

	service := NewService(mockRepo, 2, mockImageKit, mockSSEHub)

	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     1024,
	}

	job := FileUploadJob{
		Key:        "fileKey",
		FileHeader: fileHeader,
	}

	mockFile := new(MockFile)
	mockFile.On("Open").Return(nil, errors.New("failed to open file"))

	ctx := context.Background()
	result := service.processFile(ctx, job, UploadProviderImageKit)

	assert.Error(t, result.Error)
	assert.Contains(t, result.Error.Error(), "failed to open file")
	assert.Equal(t, "test.txt", result.FileName)
}

func TestProcessFile_UploadFailed(t *testing.T) {
	mockRepo := new(MockFileTransferRepo)
	mockImageKit := new(MockImageKitWrapper)
	mockSSEHub := new(sse.SSEHub)

	service := NewService(mockRepo, 2, mockImageKit, mockSSEHub)

	fileContent := "test file content"
	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     int64(len(fileContent)),
	}

	job := FileUploadJob{
		Key:        "fileKey",
		FileHeader: fileHeader,
	}

	config.Set("multipart.uploadProvider", UploadProviderImageKit)

	mockImageKit.On("UploadFile", mock.Anything, mock.Anything, "fileKey", "test.txt").Return(imagekit.UploadResponse{}, errors.New("upload failed"))

	ctx := context.Background()
	result := service.processFile(ctx, job, UploadProviderImageKit)

	assert.Error(t, result.Error)
	assert.Equal(t, "test.txt", result.FileName)
}
