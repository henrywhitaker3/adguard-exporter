package http

import (
	"context"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Http struct {
	e *echo.Echo
}

func NewHttp() *Http {
	e := echo.New()
	e.HideBanner = true

	e.GET("/metrics", echoprometheus.NewHandler())

	return &Http{
		e: e,
	}
}

func (h *Http) Serve() error {
	return h.e.Start(":9618")
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
