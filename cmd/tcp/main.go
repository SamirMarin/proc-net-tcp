package main

import (
	"proc-net-tcp/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/health", handlers.Health)

	go handlers.Tcp()

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}
