package handlers

import (
	"auth-service/internal/config"
	"auth-service/internal/tokens"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Получение токена для SDK
// @Description Получается токен для работы SDK
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} sdkTokenResponse
// @Failure 401 {object} map[string]string
// @Router /sdk-token [post]
func getSdkToken(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, ok := c.Get("user_id").(string)
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user_id is invalid"})
		}

		sdkAccessToken, err := tokens.CreateAccessToken(userId, cfg.AccessTokenSDKTTL, cfg.JWTSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed create access"})
		}
		return c.JSON(http.StatusOK, sdkTokenResponse{sdkAccessToken})
	}
}
