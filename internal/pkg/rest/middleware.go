package rest

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hafiedh.com/downloader/internal/infrastructure/container"
	utility "hafiedh.com/downloader/internal/pkg/utils"
)

type (
	DataValidator struct {
		ValidatorData *validator.Validate
	}
)

func SetupMiddleware(server *echo.Echo, container *container.Container) {
	server.Use(SetDefaultLogger())

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS, echo.PATCH, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderXCSRFToken},
		ExposeHeaders:    []string{echo.HeaderContentLength, echo.HeaderAccessControlAllowOrigin, echo.HeaderContentDisposition},
		AllowCredentials: true,
	}))

	server.HTTPErrorHandler = server.DefaultHTTPErrorHandler
	v := validator.New()
	v.RegisterValidation("ISO8601Date", utility.IsISO8601Date)
	server.Validator = &DataValidator{ValidatorData: v}
}

func SetDefaultLogger() echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/" {
				return true
			}
			if c.Request().URL.Path == "/v1/files/upload" {
				return true
			}
			return false
		},
		Handler: func(c echo.Context, reqBody []byte, resBody []byte) {
			slog.InfoContext(c.Request().Context(), "Request", "Body", string(reqBody))
			slog.InfoContext(c.Request().Context(), "Response", "Body", string(resBody))
		},
	})

}

func (cv *DataValidator) Validate(i interface{}) error {
	return cv.ValidatorData.Struct(i)
}
