package http

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Http struct {
	e *echo.Echo
}

func NewHttp(debug bool) *Http {
	e := echo.New()
	e.HideBanner = true

	e.GET("/metrics", echoprometheus.NewHandler())
	if debug {
		e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	}

	return &Http{
		e: e,
	}
}

func (h *Http) Serve() error {
	log.Println("Starting http server on port 9618")
	return h.e.Start(":9618")
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
