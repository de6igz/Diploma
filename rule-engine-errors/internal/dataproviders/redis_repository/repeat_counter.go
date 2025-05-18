package redis_repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"rule-engine-errors/internal/domain"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisRepeatCounter struct {
	rdb              *redis.Client
	maxTimeWindowSec int
	logger           *zerolog.Logger
}

func NewRedisRepeatCounter(rdb *redis.Client, maxSec int, logger *zerolog.Logger) *RedisRepeatCounter {
	return &RedisRepeatCounter{
		rdb:              rdb,
		maxTimeWindowSec: maxSec,
		logger:           logger,
	}
}

// UpdateKey – добавляет событие в ZSet, чистит старое
//func (rc *RedisRepeatCounter) UpdateKey(ctx context.Context, e *domain.Event) error {
//	rc.logger.Debug().Msgf("UpdateKey in Redis for event: %s|%s|%s|%s",
//		e.ServiceName, e.Environment, e.EventType, e.EventMessage)
//	key := rc.makeKey(e)
//	now := float64(time.Now().Unix())
//	member := fmt.Sprintf("%f-%d", now, time.Now().UnixNano())
//
//	if err := rc.rdb.ZAdd(ctx, key, redis.Z{Score: now, Member: member}).Err(); err != nil {
//		rc.logger.Error().Err(err).Msg("Failed ZAdd in Redis")
//		return err
//	}
//	cutoff := now - float64(rc.maxTimeWindowSec)
//	err := rc.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("(%f", cutoff)).Err()
//	if err != nil {
//		rc.logger.Error().Err(err).Msg("Failed ZRemRangeByScore in Redis")
//	}
//	return err
//}

//Есть ключ
// 1) Проверяем есть ли ключ
// a) Если нет то создаем со значением 1 и TTL min и возвращаем 1
// b) Если есть то INC 1 и возвращаем значение

// CountInWindow – zcount
func (rc *RedisRepeatCounter) CountInWindow(ctx context.Context, e *domain.Event, r domain.Rule, minutes int) (int, error) {
	key := rc.makeKey(e, r.ID)

	result, err := rc.rdb.Exists(ctx, key).Result()
	if err != nil {
		rc.logger.Error().Err(err).Msg("Failed Check key Exists in Redis")
		return 0, err
	}

	if result == 0 {
		rc.logger.Debug().Msgf("Key %s not exists, creating", key)
		if err = rc.rdb.Set(ctx, key, 1, time.Duration(minutes)*time.Minute).Err(); err != nil {
			rc.logger.Error().Err(err).Msg("Failed Set key in Redis")
		}
		return 1, nil
	} else {
		rc.logger.Debug().Msgf("Key %s exists, incrementing", key)
		if err = rc.rdb.Incr(ctx, key).Err(); err != nil {
			rc.logger.Error().Err(err).Msg("Failed Incr key in Redis")
		}
	}

	cnt := rc.rdb.Get(ctx, key).Val()

	rc.logger.Debug().Msgf("CountInWindow for key=%s, minutes=%d", key, minutes)
	cntInt, err := strconv.Atoi(cnt)
	if err != nil {
		rc.logger.Error().Err(err).Msg("Failed to convert cnt to int")
		return 0, err
	}

	return cntInt, nil

}

// makeKey
// cache-key: rule_id-user_id-service_name_environment
func (rc *RedisRepeatCounter) makeKey(e *domain.Event, ruleId string) string {
	return fmt.Sprintf("%s-%s-%s-%s", ruleId, e.UserID, e.ServiceName, e.Environment)
}
