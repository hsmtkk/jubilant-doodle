package main

import (
	"fmt"
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
	return c.String(http.StatusOK, "pong")
}
