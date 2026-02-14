package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microservices-boilerplate/go/internal/config"
)

// HealthHandler provides health check endpoints.
type HealthHandler struct {
	cfg *config.Config
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

// Liveness returns a simple alive check.
func (h *HealthHandler) Liveness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "alive",
		"service": h.cfg.AppName,
	})
}

// Readiness returns a readiness check (extend for DB, etc.).
func (h *HealthHandler) Readiness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ready",
		"service": h.cfg.AppName,
	})
}
