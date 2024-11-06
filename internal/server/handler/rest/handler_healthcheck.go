package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hafiedh.com/downloader/internal/usecase/healthcheck"
)

type healthCheckHandler struct {
	healthCheckService healthcheck.Service
}

func NewHealthCheckHandler() *healthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) SetHealthCheckService(service healthcheck.Service) *healthCheckHandler {
	h.healthCheckService = service
	return h
}

func (h *healthCheckHandler) Validate() *healthCheckHandler {
	if h.healthCheckService == nil {
		panic("healthCheckService is nil")
	}
	return h
}

func (h *healthCheckHandler) HealthCheck(ec echo.Context) (err error) {
	ctx := ec.Request().Context()

	res, err := h.healthCheckService.HealthCheck(ctx)
	if err != nil {
		return
	}

	return ec.JSON(http.StatusOK, res)
}
