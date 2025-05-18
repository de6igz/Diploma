package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// CreateAccessToken создаёт JWT с заданным TTL
func CreateAccessToken(userID string, ttl time.Duration, secret string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken парсит и валидирует JWT
func ValidateAccessToken(tokenStr, secret string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("invalid subject")
		}
		return sub, nil
	}
	return "", fmt.Errorf("invalid token")
}
