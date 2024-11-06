package container

import (
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/keepalive"
	"hafiedh.com/downloader/internal/config"
	"hafiedh.com/downloader/internal/domain/repositories"
	"hafiedh.com/downloader/internal/infrastructure/imagekit"
	"hafiedh.com/downloader/internal/infrastructure/postgresql"
	"hafiedh.com/downloader/internal/infrastructure/sse"
	"hafiedh.com/downloader/internal/usecase/downloader"
	"hafiedh.com/downloader/internal/usecase/healthcheck"
)

type Container struct {
	Config             *config.DefaultConfig
	PostgresqlDB       *config.PostgreSQLDB
	HealthCheckService healthcheck.Service
	DownloaderService  downloader.FileTransfer
	SSEHub             *sse.SSEHub
	Logger             *slog.Logger
	Echo               *echo.Echo
	GrpcConfig         *config.GrpcServerConfig
	PsqlDB             *postgresql.PostgresImpl
}

func (c *Container) Validate() *Container {
	if c.Config == nil {
		panic("Config is nil")
	}
	if c.HealthCheckService == nil {
		panic("HealthCheckService is nil")
	}
	if c.Logger == nil {
		panic("Logger is nil")
	}
	if c.Echo == nil {
		panic("Echo is nil")
	}
	if c.GrpcConfig == nil {
		panic("GrpcServer is nil")
	}
	if c.PsqlDB == nil {
		panic("PostgresqlDB is nil")
	}
	return c
}

func New() *Container {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	config.Load(os.Getenv("env"), ".env")
	echo := echo.New()

	_ = config.GetString("logger.fileLocation")
	_ = config.GetString("logger.fileTdrLocation")
	_ = time.Duration(config.GetInt("logger.fileMaxAge"))
	_ = config.GetBool("logger.stdout")

	defConfig := &config.DefaultConfig{
		Apps: config.Apps{
			Name:     config.GetString("app.name"),
			Address:  config.GetString("address"),
			HttpPort: config.GetString("port"),
		},
	}

	grpcConfig := &config.GrpcServerConfig{
		Port: config.GetInt("grpc.port"),
		KeepaliveParams: keepalive.ServerParameters{
			MaxConnectionIdle:     100,
			MaxConnectionAge:      7200,
			MaxConnectionAgeGrace: 60,
			Time:                  10,
			Timeout:               3,
		},

		KeepalivePolicy: keepalive.EnforcementPolicy{
			MinTime:             10,
			PermitWithoutStream: true,
		},
	}

	pgConfig := &config.PostgreSQLDB{
		Username:                 config.GetString("postgresql.downloader.username"),
		Password:                 config.GetString("postgresql.downloader.password"),
		Name:                     config.GetString("postgresql.downloader.db"),
		Schema:                   config.GetString("postgresql.downloader.schema"),
		Host:                     config.GetString("postgresql.downloader.host"),
		Port:                     config.GetInt("postgresql.downloader.port"),
		DefaultMaxConn:           config.GetInt("postgresql.downloader.maxConns"),
		DefaultMinConn:           config.GetInt("postgresql.downloader.minConns"),
		DefaultConnectTimeout:    time.Duration(config.GetInt("postgresql.downloader.timeOut") * int(time.Second)),
		DefaultMaxConnLifetime:   time.Hour,
		DefaultMaxConnIdleTime:   30 * time.Minute,
		DefaultHealthCheckPeriod: time.Minute,
	}

	db, err := postgresql.NewConnection(pgConfig)
	if err != nil {
		slog.Error("Failed to create connection pool", "error", err)
		panic(err)
	}

	// * Repositories
	downloaderRepo := repositories.NewFileTransfer(db.Pool)

	// * Wrapper
	imageKit := imagekit.NewImageKitWrapper()
	sseHub := sse.NewSSEHub()

	// * Services
	healthCheckService := healthcheck.NewService().Validate()
	downloaderService := downloader.NewService(downloaderRepo, 10, imageKit, sseHub)

	// * Brokers

	// * Workers

	container := &Container{
		Config:             defConfig,
		HealthCheckService: healthCheckService,
		DownloaderService:  downloaderService,
		Logger:             logger,
		Echo:               echo,
		GrpcConfig:         grpcConfig,
		PostgresqlDB:       pgConfig,
		PsqlDB:             &db,
		SSEHub:             sseHub,
	}
	container.Validate()
	return container

}
