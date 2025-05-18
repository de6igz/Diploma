package handlers

import (
	"auth-service/internal/config"
	"auth-service/internal/tokens"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Валидация токена (Nginx Auth)
// @Description Проверяет валидность access токена. Возвращает X-User-Id при успехе
// @Tags Auth
// @Success 200
// @Failure 401
// @Router /validate [get]
func validateHandler(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		tokenStr := extractBearerToken(authHeader)
		if tokenStr == "" {
			return c.String(http.StatusUnauthorized, "missing or invalid Authorization header")
		}
		userID, err := tokens.ValidateAccessToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			return c.String(http.StatusUnauthorized, "invalid token")
		}
		c.Response().Header().Set("X-User-Id", userID)
		return c.NoContent(http.StatusOK)
	}
}
