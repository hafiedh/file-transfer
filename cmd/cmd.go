package cmd

import (
	"hafiedh.com/downloader/internal/infrastructure/container"
	"hafiedh.com/downloader/internal/server"
)

func Run() {
	server.StartService(container.New())
}
