package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Health
// @Description health
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200
// @Router /health [get]
func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
