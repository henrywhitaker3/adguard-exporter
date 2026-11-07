package http

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"

	echoprometheus "github.com/labstack/echo-prometheus"
	"github.com/labstack/echo/v5"
)

type HTTP struct {
	e *echo.Echo

	ready    bool
	healthy  bool
	healthMu *sync.Mutex
}

func NewHttp(debug bool) *HTTP {
	e := echo.New()

	e.GET("/metrics", echoprometheus.NewHandler())
	if debug {
		e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	}

	http := &HTTP{
		e:        e,
		ready:    false,
		healthMu: &sync.Mutex{},
	}

	http.e.GET("/healthz", http.healthz())
	http.e.GET("/readyz", http.readyz())

	return http
}

func (h *HTTP) Serve(bindAddr string) error {
	log.Println("Starting http server on " + bindAddr)
	return h.e.Start(bindAddr)
}

func (h *HTTP) Ready(state bool) {
	h.healthMu.Lock()
	defer h.healthMu.Unlock()
	h.ready = state
}

func (h *HTTP) Healthy(state bool) {
	h.healthMu.Lock()
	defer h.healthMu.Unlock()
	h.healthy = state
}

func (h *HTTP) healthz() echo.HandlerFunc {
	return func(c *echo.Context) error {
		h.healthMu.Lock()
		defer h.healthMu.Unlock()
		code := http.StatusOK
		if !h.healthy {
			code = http.StatusServiceUnavailable
		}
		return c.NoContent(code)
	}
}

func (h *HTTP) readyz() echo.HandlerFunc {
	return func(c *echo.Context) error {
		h.healthMu.Lock()
		defer h.healthMu.Unlock()
		code := http.StatusOK
		if !h.ready {
			code = http.StatusServiceUnavailable
		}
		return c.NoContent(code)
	}
}
