package server

import (
	"fmt"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/gommon/color"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"hafiedh.com/downloader/internal/infrastructure/container"
	rest "hafiedh.com/downloader/internal/pkg/rest"
	router "hafiedh.com/downloader/internal/server/handler/rest"
)

func StartHttpServer(container *container.Container) {
	rest.SetupMiddleware(container.Echo, container)
	router.SetupHandler(container)
	router.SetupRouter(container.Echo, container)
	go container.SSEHub.Run()

	container.Echo.Server.Addr = fmt.Sprintf("%s:%s", container.Config.Apps.Address, container.Config.Apps.HttpPort)

	color.Println(color.Green(fmt.Sprintf("⇨ http started on port: %s\n", container.Config.Apps.HttpPort)))

	err := gracehttp.Serve(&http.Server{Addr: container.Echo.Server.Addr, Handler: h2c.NewHandler(container.Echo, &http2.Server{MaxConcurrentStreams: 500, MaxReadFrameSize: 1048576})})
	if err != nil {
		panic(err)
	}
}
