package http

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Http struct {
	e *echo.Echo

	ready    bool
	healthy  bool
	healthMu *sync.Mutex
}

func NewHttp(debug bool) *Http {
	e := echo.New()
	e.HideBanner = true

	e.GET("/metrics", echoprometheus.NewHandler())
	if debug {
		e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	}

	http := &Http{
		e:        e,
		ready:    false,
		healthMu: &sync.Mutex{},
	}

	http.e.GET("/healthz", http.healthz())
	http.e.GET("/readyz", http.readyz())

	return http
}

func (h *Http) Serve(bindAddr string) error {
	log.Println("Starting http server on " + bindAddr)
	return h.e.Start(bindAddr)
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}

func (h *Http) Ready(state bool) {
	h.healthMu.Lock()
	defer h.healthMu.Unlock()
	h.ready = state
}

func (h *Http) Healthy(state bool) {
	h.healthMu.Lock()
	defer h.healthMu.Unlock()
	h.healthy = state
}

func (h *Http) healthz() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.healthMu.Lock()
		defer h.healthMu.Unlock()
		code := http.StatusOK
		if !h.healthy {
			code = http.StatusServiceUnavailable
		}
		return c.NoContent(code)
	}
}

func (h *Http) readyz() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.healthMu.Lock()
		defer h.healthMu.Unlock()
		code := http.StatusOK
		if !h.ready {
			code = http.StatusServiceUnavailable
		}
		return c.NoContent(code)
	}
}
