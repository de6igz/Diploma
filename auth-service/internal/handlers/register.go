package handlers

import (
	"auth-service/internal/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Регистрация
// @Description Регистрирует нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body registerRequest true "Регистрационные данные"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /register [post]
func registerHandler(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	var userID int64
	err := db.DBConn.QueryRow(query, req.Username, req.Password).Scan(&userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "registered",
		"user_id": userID,
	})
}
