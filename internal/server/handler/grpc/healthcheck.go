package grpc

import (
	"context"
	"log/slog"

	pb "hafiedh.com/downloader/gen/go/healthcheck/v2"
	"hafiedh.com/downloader/internal/usecase/healthcheck"
)

func NewHealthCheckHandler(healthCheckService healthcheck.Service) pb.HealthCheckServiceServer {
	if healthCheckService == nil {
		panic("healthCheckService is nil")
	}
	return &healthCheckHandler{
		healthCheckService: healthCheckService,
	}
}

func (h *healthCheckHandler) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	res, err := h.healthCheckService.HealthCheck(ctx)
	if err != nil {
		slog.Error("failed to health check err=%v\n", err)
		return nil, err
	}

	return &pb.HealthCheckResponse{
		Message:    res.Message,
		Version:    res.Version,
		ServerTime: res.ServerTime,
	}, nil

}
