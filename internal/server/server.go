package server

import (
	"io"
	"log/slog"

	"os"
	"os/signal"
	"syscall"

	googleGrpc "google.golang.org/grpc"
	pbFile "hafiedh.com/downloader/gen/go/file_transfer/v2"
	pb "hafiedh.com/downloader/gen/go/healthcheck/v2"
	"hafiedh.com/downloader/internal/infrastructure/container"
)

func StartService(container *container.Container) {
	grpcServer, handler, err := NewGrpcServer(container)
	if err != nil {
		slog.Error("failed to create grpc server err=%v\n", err)
		return
	}

	go StartHttpServer(container)

	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			pb.RegisterHealthCheckServiceServer(server, handler.HealthCheckHandler)
			pbFile.RegisterFileTransferServiceServer(server, handler.FileTransferHandler)
		},
	)
	AddShutdownHook(grpcServer, container.Echo, container.PsqlDB)
}

func AddShutdownHook(closers ...io.Closer) {
	slog.Info("registering shutdown hook", "closers", closers)
	c := make(chan os.Signal, 1)
	signal.Notify(
		c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)

	<-c
	slog.Info("received signal, starting graceful shutdown")

	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			slog.Error("failed to close resource err=%v\n", err)
		}
	}
	slog.Info("shutdown completed")
}
