package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"

	"auth-service/internal/config"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
	ctx := context.Background()
	return RedisClient.Ping(ctx).Err()
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
