package grpc

import (
	pb "hafiedh.com/downloader/gen/go/healthcheck/v2"
	"hafiedh.com/downloader/internal/infrastructure/container"
	"hafiedh.com/downloader/internal/usecase/healthcheck"
)

type (
	healthCheckHandler struct {
		pb.UnimplementedHealthCheckServiceServer
		healthCheckService healthcheck.Service
	}

	Handler struct {
		HealthCheckHandler *healthCheckHandler
	}
)

func SetupHandler(container *container.Container) Handler {
	return Handler{
		HealthCheckHandler: NewHealthCheckHandler(container.HealthCheckService).(*healthCheckHandler),
	}
}
