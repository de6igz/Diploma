package utils

import (
	"auth-service/internal/db"
	"auth-service/internal/models"
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"
)

// ====================== Вспомогательные ======================

func getUserByUsername(username string) (*models.User, error) {
	var u models.User
	query := `SELECT id, username, password FROM users WHERE username=$1 LIMIT 1`
	err := db.DBConn.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password)
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
