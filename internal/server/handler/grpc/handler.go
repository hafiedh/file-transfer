package grpc

import (
	pbFile "hafiedh.com/downloader/gen/go/file_transfer/v2"
	pb "hafiedh.com/downloader/gen/go/healthcheck/v2"
	"hafiedh.com/downloader/internal/infrastructure/container"
	"hafiedh.com/downloader/internal/usecase/downloader"
	"hafiedh.com/downloader/internal/usecase/healthcheck"
)

type (
	healthCheckHandler struct {
		pb.UnimplementedHealthCheckServiceServer
		healthCheckService healthcheck.Service
	}

	fileTransferHandler struct {
		pbFile.UnimplementedFileTransferServiceServer
		fileTransferService downloader.FileTransfer
	}

	Handler struct {
		HealthCheckHandler  *healthCheckHandler
		FileTransferHandler *fileTransferHandler
	}
)

func SetupHandler(container *container.Container) Handler {
	return Handler{
		HealthCheckHandler:  NewHealthCheckHandler(container.HealthCheckService).(*healthCheckHandler),
		FileTransferHandler: NewFileTransferHandler(container.DownloaderService).(*fileTransferHandler),
	}
}
