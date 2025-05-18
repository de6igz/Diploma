package redis_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rule-engine-errors/internal/domain"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// RedisCache отвечает за кеширование правил в Redis.
type RedisCache struct {
	rdb        *redis.Client
	defaultTTL time.Duration
	logger     *zerolog.Logger
}

// NewRedisRuleCache создаёт новый экземпляр кеша правил.
func NewRedisRuleCache(rdb *redis.Client, defaultTTL time.Duration, logger *zerolog.Logger) *RedisCache {
	return &RedisCache{
		rdb:        rdb,
		defaultTTL: defaultTTL,
		logger:     logger,
	}
}

// GetRules пытается получить правила из кеша по комбинации userID и serviceName.
// Если ключ не найден, возвращается nil (и можно будет далее обращаться к Mongo).
func (rc *RedisCache) GetRules(ctx context.Context, userID, serviceName, projectId string) ([]domain.Rule, error) {
	key := rc.makeKey(userID, serviceName, projectId)
	data, err := rc.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		rc.logger.Debug().Msgf("Cache miss for key %s", key)
		return nil, nil // Ключ не найден – кеш промах.
	} else if err != nil {
		rc.logger.Error().Err(err).Msgf("Failed to get key %s from Redis", key)
		return nil, err
	}

	var rules []domain.Rule
	if err := json.Unmarshal([]byte(data), &rules); err != nil {
		rc.logger.Error().Err(err).Msgf("Failed to unmarshal rules for key %s", key)
		return nil, err
	}
	rc.logger.Debug().Msgf("Cache hit for key %s", key)
	return rules, nil
}

// SetRules записывает правила в кеш с заданным TTL.
func (rc *RedisCache) SetRules(ctx context.Context, userID, serviceName, projectId string, rules []domain.Rule) error {
	key := rc.makeKey(userID, serviceName, projectId)
	data, err := json.Marshal(rules)
	if err != nil {
		rc.logger.Error().Err(err).Msg("Failed to marshal rules for cache")
		return err
	}
	if err := rc.rdb.Set(ctx, key, data, rc.defaultTTL*time.Second).Err(); err != nil {
		rc.logger.Error().Err(err).Msgf("Failed to set cache for key %s", key)
		return err
	}
	rc.logger.Debug().Msgf("Cache set for key %s with TTL %s", key, rc.defaultTTL)
	return nil
}

// makeKey формирует ключ кеша для комбинации userID и serviceName.
func (rc *RedisCache) makeKey(userID, serviceName, projectId string) string {
	return fmt.Sprintf("rules-errors:%s:%s:%s", userID, serviceName, projectId)
}
