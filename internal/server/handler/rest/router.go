package rest

import (
	"github.com/labstack/echo/v4"
	"hafiedh.com/downloader/internal/infrastructure/container"
)

func SetupRouter(e *echo.Echo, cnt *container.Container) {
	h := SetupHandler(cnt).Validate()

	e.GET("/", h.healthCheckHandler.HealthCheck)

	e.GET("/sse", h.sseHandler.SSE)

	base := e.Group("/v1")

	files := base.Group("/files")
	{
		files.POST("/upload", h.downloaderHandler.Uploader)
		files.GET("/:id/presigned", h.downloaderHandler.Presigned)
	}

}
