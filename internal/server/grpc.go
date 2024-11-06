package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/labstack/gommon/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"hafiedh.com/downloader/internal/config"
	"hafiedh.com/downloader/internal/infrastructure/container"
	grpcHandler "hafiedh.com/downloader/internal/server/handler/grpc"
)

type (
	GRPCServer interface {
		Start(serviceRegister func(server *grpc.Server))
		io.Closer
	}

	gRPCServer struct {
		grpcServer *grpc.Server
		config     config.GrpcServerConfig
	}
)

func NewGrpcServer(container *container.Container) (server GRPCServer, handler grpcHandler.Handler, err error) {
	options, err := buildOptions(*container.GrpcConfig)
	if err != nil {
		return
	}

	grpcServer := grpc.NewServer(options...)
	return &gRPCServer{grpcServer: grpcServer, config: *container.GrpcConfig}, grpcHandler.SetupHandler(container), nil
}

func (s *gRPCServer) Start(serviceRegister func(server *grpc.Server)) {
	grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(int(s.config.Port)))
	if err != nil {
		panic(err)
	}

	serviceRegister(s.grpcServer)

	color.Println(color.Green(fmt.Sprintf("⇨ gRPC started on port: %d\n", s.config.Port)))
	if err := s.grpcServer.Serve(grpcListener); err != nil {
		panic(err)
	}
}

func (s *gRPCServer) Close() error {
	s.grpcServer.GracefulStop()
	return nil
}

func buildOptions(config config.GrpcServerConfig) ([]grpc.ServerOption, error) {
	return []grpc.ServerOption{
		grpc.KeepaliveParams(buildKeepaliveParams(config.KeepaliveParams)),
		grpc.KeepaliveEnforcementPolicy(buildKeepalivePolicy(config.KeepalivePolicy)),
	}, nil
}

func buildKeepalivePolicy(config keepalive.EnforcementPolicy) keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             config.MinTime * time.Second,
		PermitWithoutStream: config.PermitWithoutStream,
	}
}

func buildKeepaliveParams(config keepalive.ServerParameters) keepalive.ServerParameters {
	return keepalive.ServerParameters{
		MaxConnectionIdle:     config.MaxConnectionIdle * time.Second,
		MaxConnectionAge:      config.MaxConnectionAge * time.Second,
		MaxConnectionAgeGrace: config.MaxConnectionAgeGrace * time.Second,
		Time:                  config.Time * time.Second,
		Timeout:               config.Timeout * time.Second,
	}
}
