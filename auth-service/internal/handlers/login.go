package handlers

import (
	"auth-service/internal/config"
	"auth-service/internal/tokens"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Логин
// @Description Логин пользователя и выдача access/refresh токенов
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Логин данные"
// @Success 200 {object} loginResponse
// @Failure 401 {object} map[string]string
// @Router /login [post]
func loginHandler(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req loginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		user, err := getUserByUsername(req.Username)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}
		if user.Password != req.Password {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}

		// Определяем TTL для access

		ttl := cfg.AccessTokenFrontTTL

		accessToken, err := tokens.CreateAccessToken(fmt.Sprintf("%d", user.ID), ttl, cfg.JWTSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed create access"})
		}
		refreshToken, err := createRefreshToken()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed create refresh"})
		}
		if err := storeRefreshToken(refreshToken, fmt.Sprintf("%d", user.ID), cfg.RefreshTokenTTL); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed store refresh"})
		}

		return c.JSON(http.StatusOK, loginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
