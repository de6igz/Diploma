package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	// PostgreSQL
	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDB       string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// JWT
	JWTSecret string

	// TTL
	AccessTokenFrontTTL time.Duration
	AccessTokenSDKTTL   time.Duration
	RefreshTokenTTL     time.Duration

	// Service
	ServicePort string
}

// LoadConfig загружает конфигурацию из переменных окружения.
// Если какая-либо обязательная переменная не задана или имеет неверный формат, программа завершится с ошибкой.
func LoadConfig() *Config {
	return &Config{
		// PostgreSQL
		PGHost:     mustGetEnv("PGHOST"),
		PGPort:     mustGetEnv("PGPORT"),
		PGUser:     mustGetEnv("PGUSER"),
		PGPassword: mustGetEnv("PGPASSWORD"),
		PGDB:       mustGetEnv("PGDB"),

		// Redis
		RedisHost:     mustGetEnv("REDISHOST"),
		RedisPort:     mustGetEnv("REDISPORT"),
		RedisPassword: mustGetEnv("REDISPASSWORD"),

		// JWT
		JWTSecret: mustGetEnv("JWTSECRET"),

		// TTL
		AccessTokenFrontTTL: mustParseDuration("ACCESS_TOKEN_FRONT_TTL"),
		AccessTokenSDKTTL:   mustParseDuration("ACCESS_TOKEN_SDK_TTL"),
		RefreshTokenTTL:     mustParseDuration("REFRESH_TOKEN_TTL"),

		// Service
		ServicePort: mustGetEnv("PORT"),
	}
}

// mustGetEnv получает значение переменной окружения по ключу.
// Если переменная не задана, выводит ошибку и завершает программу.
func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Ошибка: переменная окружения %s не задана", key)
	}
	return val
}

// mustParseDuration получает значение переменной окружения и преобразует его в time.Duration.
// Если переменная не задана или имеет неверный формат, выводит ошибку и завершает программу.
func mustParseDuration(key string) time.Duration {
	val := mustGetEnv(key)
	duration, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("Ошибка: неверный формат для %s: %v", key, err)
	}
	return duration
}
