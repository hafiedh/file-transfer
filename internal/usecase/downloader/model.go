package downloader

import (
	"io"
	"mime/multipart"
)

type (
	UploadProgress struct {
		FileName   string  `json:"fileName"`
		Status     string  `json:"status"`
		URL        string  `json:"url,omitempty"`
		Error      string  `json:"error,omitempty"`
		Percentage float64 `json:"percentage"`
	}

	ProgressReader struct {
		reader     io.Reader
		onProgress func(UploadProgress)
		total      int64
		read       int64
	}

	FileUploadJob struct {
		Key        string
		FileHeader *multipart.FileHeader
		Result     chan UploadResult
	}

	UploadResult struct {
		FileName string
		URL      string
		Error    error
	}
)

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.read += int64(n)

	progress := UploadProgress{
		Percentage: float64(pr.read) / float64(pr.total) * 100,
	}
	pr.onProgress(progress)

	return
}

func NewProgressReader(r io.Reader, total int64, onProgress func(UploadProgress)) *ProgressReader {
	return &ProgressReader{
		reader:     r,
		total:      total,
		onProgress: onProgress,
	}
}
