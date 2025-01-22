package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthStatus struct {
	Health bool `json:"health"`
}

func (h *Handler) HealthCheck(c echo.Context) error {
	//Anonymous struct -> check health
	message := HealthStatus{
		Health: true,
	}

	return c.JSON(http.StatusOK, message)
}
