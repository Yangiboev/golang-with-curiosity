package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *server) runHttpServer() {
	s.echo.GET("health", func(c *echo.Context) {
		c.String(http.StatusOK, "Ok")
	})
	s.echo.GET("/metrics", echo.WrapH(promhttp.Handler()))
	s.mapRoutes()
}
func (s *server) mapRoutes() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Products microservice"
	docs.SwaggerInfo.Description = "Products REST API microservice."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.Use(echo.Logger())
	s.echo.Use(echo.HTTPSRedirect())
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrfTokenHeader},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit(bodyLimit))
}
