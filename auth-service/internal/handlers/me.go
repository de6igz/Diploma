package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Получить инфу о себе
// @Description Получить инфу о себе
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200
// @Router /me [get]
func meHandler(c echo.Context) error {
	userId := c.Get("user_id").(string)
	user, err := getUserById(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       userId,
		"userName": user.Username,
	})

}
