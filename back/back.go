package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/hsmtkk/jubilant-doodle/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port, err := env.GetPort()
	if err != nil {
		log.Fatal(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", ping)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

// Handler
func ping(c echo.Context) error {
	reqBytes, err := httputil.DumpRequest(c.Request(), false)
	if err != nil {
		return fmt.Errorf("httputil.DumpRequest failed; %w", err)
	}
	fmt.Println(string(reqBytes))
	return c.String(http.StatusOK, "pong")
}
