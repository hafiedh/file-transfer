package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"hafiedh.com/downloader/internal/infrastructure/sse"
	"hafiedh.com/downloader/internal/pkg/constants"
)

type (
	sseHandler struct {
		sseHub *sse.SSEHub
	}
)

func NewSSEHandler(sseHub *sse.SSEHub) *sseHandler {
	return &sseHandler{
		sseHub: sseHub,
	}
}

func (s *sseHandler) SSE(ec echo.Context) (err error) {
	client := sse.SSEClient{
		ID:      ec.Response().Header().Get(echo.HeaderXRequestID),
		Message: make(chan constants.DefaultResponse),
	}

	s.sseHub.AddClient(&client)
	defer func() { s.sseHub.RemoveClient(&client) }()

	ec.Response().Header().Set("Content-Type", "text/event-stream")
	ec.Response().Header().Set("Cache-Control", "no-cache")
	ec.Response().Header().Set("Connection", "keep-alive")

	if _, ok := ec.Response().Writer.(http.Flusher); !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")
	}

	for msg := range client.Message {
		data, err := json.Marshal(msg)
		if err != nil {
			slog.Error("failed to marshal message", "error", err)
			continue
		}
		fmt.Fprintf(ec.Response().Writer, "%s", data)
		fmt.Fprintf(ec.Response().Writer, "\n\n")
		if f, ok := ec.Response().Writer.(http.Flusher); ok {
			f.Flush()
		}
	}

	return
}
