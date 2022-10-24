package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hsmtkk/jubilant-doodle/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port, err := env.GetPort()
	if err != nil {
		log.Fatal(err)
	}

	dstURL, err := env.RequiredEnv("DST_URL")
	if err != nil {
		log.Fatal(err)
	}

	hdl := newHandler(dstURL)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hdl.ping)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type myHandler struct {
	dstURL string
}

func newHandler(dstURL string) *myHandler {
	return &myHandler{dstURL}
}

// Handler
func (h *myHandler) ping(c echo.Context) error {
	resp, err := http.Get(h.dstURL)
	if err != nil {
		return fmt.Errorf("http.Get failed; %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed; %w", err)
	}

	return c.String(http.StatusOK, string(respBytes))
}
