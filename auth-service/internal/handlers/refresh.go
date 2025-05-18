package handlers

import (
	"auth-service/internal/config"
	"auth-service/internal/tokens"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Обновление токена
// @Description Обновляет access/refresh токены
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshRequest body refreshRequest true "Refresh данные"
// @Success 200 {object} refreshResponse
// @Failure 401 {object} map[string]string
// @Router /refresh [post]
func refreshHandler(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req refreshRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		userID, err := getUserIDByRefreshToken(req.RefreshToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh"})
		}

		err = deleteRefreshToken(req.RefreshToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed delete refresh token"})
		}

		// Генерируем новый refresh
		newRefresh, err := createRefreshToken()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create refresh"})
		}
		if err := storeRefreshToken(newRefresh, userID, cfg.RefreshTokenTTL); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed store refresh"})
		}

		// Генерируем новый access
		ttl := cfg.AccessTokenFrontTTL

		accessToken, err := tokens.CreateAccessToken(userID, ttl, cfg.JWTSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed create access"})
		}

		return c.JSON(http.StatusOK, refreshResponse{
			AccessToken:  accessToken,
			RefreshToken: newRefresh,
		})
	}
}
