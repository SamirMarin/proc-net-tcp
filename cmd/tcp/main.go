package main

import (
	"proc-net-tcp/handlers"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	go handlers.Tcp()

	e := echo.New()
	e.GET("/health", handlers.Health)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}
