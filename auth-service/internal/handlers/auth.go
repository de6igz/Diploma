package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/models"
	"auth-service/internal/tokens"
)

// RegisterAuthRoutes регистрирует маршруты аутентификации в Echo
func RegisterAuthRoutes(e *echo.Echo, cfg *config.Config) {
	e.POST("/register", registerHandler)
	e.POST("/login", loginHandler(cfg))
	e.POST("/refresh", refreshHandler(cfg))
	e.GET("/validate", validateHandler(cfg))
	e.POST("/sdk-token", getSdkToken(cfg), withAccessToken(cfg))
	// Пример защищённого эндпоинта
	e.GET("/secure", secureEndpoint, withAccessToken(cfg))
	e.GET("/health", healthHandler)
	e.GET("/me", meHandler, withAccessToken(cfg))
}

// Пример защищённого эндпоинта
// @Summary Пример защищённого эндпоинта
// @Description Возвращает секретные данные, если токен валиден
// @Tags Auth
// @Accept json
// @Produce text/plain
// @Success 200 {string} string "Secure data"
// @Failure 401
// @Router /secure [get]
func secureEndpoint(c echo.Context) error {
	userID := c.Get("user_id").(string)
	return c.String(http.StatusOK, fmt.Sprintf("Secure Data for user %s", userID))
}

// Middleware проверяющий access токен
func withAccessToken(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			tokenStr := extractBearerToken(authHeader)
			if tokenStr == "" {
				return c.String(http.StatusUnauthorized, "missing token")
			}
			userID, err := tokens.ValidateAccessToken(tokenStr, cfg.JWTSecret)
			if err != nil {
				return c.String(http.StatusUnauthorized, "invalid token")
			}
			c.Set("user_id", userID)
			return next(c)
		}
	}
}

// ====================== Вспомогательные ======================

// Удаляет refresh токен из Redis
func deleteRefreshToken(refreshToken string) error {
	ctx := context.Background()
	return db.RedisClient.Del(ctx, refreshToken).Err()
}

func getUserByUsername(username string) (*models.User, error) {
	var u models.User
	query := `SELECT id, username, password FROM users WHERE username=$1 LIMIT 1`
	err := db.DBConn.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func getUserById(id string) (*models.User, error) {
	var u models.User
	query := `SELECT id, username FROM users WHERE id=$1 LIMIT 1`
	err := db.DBConn.QueryRow(query, id).Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func createRefreshToken() (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

func storeRefreshToken(refreshToken, userID string, ttl time.Duration) error {
	ctx := context.Background()
	return db.RedisClient.Set(ctx, refreshToken, userID, ttl).Err()
}

func getUserIDByRefreshToken(refreshToken string) (string, error) {
	ctx := context.Background()
	val, err := db.RedisClient.Get(ctx, refreshToken).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func extractBearerToken(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return parts[1]
	}
	return ""
}
