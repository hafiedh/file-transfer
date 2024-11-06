package rest

import (
	"hafiedh.com/downloader/internal/infrastructure/container"
)

type Handler struct {
	healthCheckHandler *healthCheckHandler
	downloaderHandler  *downloadHandler
	sseHandler         *sseHandler
}

func SetupHandler(container *container.Container) *Handler {
	return &Handler{
		healthCheckHandler: NewHealthCheckHandler().SetHealthCheckService(container.HealthCheckService).Validate(),
		downloaderHandler:  NewDownloaderHandler(container.DownloaderService),
		sseHandler:         NewSSEHandler(container.SSEHub),
	}
}

func (h *Handler) Validate() *Handler {
	if h.healthCheckHandler == nil {
		panic("healthCheckHandler is nil")
	}
	if h.downloaderHandler == nil {
		panic("downloaderHandler is nil")
	}
	if h.sseHandler == nil {
		panic("sseHandler is nil")
	}
	return h
}
