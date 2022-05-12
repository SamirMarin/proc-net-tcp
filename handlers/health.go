package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthCheckResponse struct {
	Message string `json:"message"`
}

// HealthCheck - Health Check Handler
func Health(c echo.Context) error {
	resp := HealthCheckResponse{
		Message: "Healthy",
	}
	return c.JSON(http.StatusOK, resp)
}
