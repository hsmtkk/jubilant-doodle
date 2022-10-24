package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"

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

	enableAuth := false
	enableAuthStr := os.Getenv("ENABLE_AUTH")
	if enableAuthStr != "" {
		enableAuth = true
	}

	hdl := newHandler(dstURL, enableAuth)

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
	dstURL     string
	enableAuth bool
}

func newHandler(dstURL string, enableAuth bool) *myHandler {
	return &myHandler{dstURL, enableAuth}
}

// https://cloud.google.com/run/docs/authenticating/service-to-service?hl=ja#acquire-token
func (h *myHandler) ping(c echo.Context) error {
	var client *http.Client = http.DefaultClient
	if h.enableAuth {
		var err error
		client, err = idtoken.NewClient(c.Request().Context(), h.dstURL)
		if err != nil {
			return fmt.Errorf("idtoken.NewClient failed; %w", err)
		}
	}

	resp, err := client.Get(h.dstURL)
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
