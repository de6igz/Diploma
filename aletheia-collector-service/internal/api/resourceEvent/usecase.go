// resourceEvent/usecase.go
package resourceEvent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/redis/go-redis/v9"
)

// UseCase – бизнес-логика для ResourceEvent.
type UseCase interface {
	ProcessResourceEvent(evt ResourceEvent, userID string) error
}

// useCaseImpl – конкретная реализация UseCase.
type useCaseImpl struct {
	repo                 Repository
	maxRequestsPerMinute int
	redisClient          *redis.Client
}

// NewUseCase – конструктор.
func NewUseCase(repo Repository, maxRequestsPerMinute int) UseCase {
	return &useCaseImpl{
		repo:                 repo,
		maxRequestsPerMinute: maxRequestsPerMinute,
		redisClient:          repo.GetRedisClient(), // Предполагаем, что репозиторий предоставляет доступ к Redis
	}
}

// ProcessResourceEvent – обрабатывает событие ResourceEvent с ограничением скорости.
func (uc *useCaseImpl) ProcessResourceEvent(evt ResourceEvent, userID string) error {
	// Проверка ограничения скорости
	if !uc.allowRequest(userID, evt.ServiceName) {
		return fmt.Errorf(rateLimitError)
	}

	// Генерируем ключ
	key, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed to generate Kafka key for resource event: %v", err)
	}

	value, err := BuildKafkaMessage(evt, userID)
	if err != nil {
		return fmt.Errorf("failed to build message for resource event: %v", err)
	}

	partition, offset, err := uc.repo.SendResourceEvent([]byte(key), value)
	if err != nil {
		return err
	}

	// Логируем или обрабатываем partition/offset
	fmt.Printf("[ResourceEvent] Sent to topic partition=%d, offset=%d\n", partition, offset)

	return nil
}

// allowRequest проверяет, допустим ли запрос для данного userID и serviceName.
func (uc *useCaseImpl) allowRequest(userID, serviceName string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:%s:%s", userID, serviceName)

	// Используем INCR и устанавливаем TTL, если ключ новый
	// Возвращаемое значение - количество запросов
	count, err := uc.redisClient.Incr(ctx, key).Result()
	if err != nil {
		// В случае ошибки позволяем запрос, но логируем
		fmt.Printf("Redis error: %v\n", err)
		return true
	}

	if count == 1 {
		// Устанавливаем TTL на 1 минуту
		err := uc.redisClient.Expire(ctx, key, time.Minute).Err()
		if err != nil {
			fmt.Printf("Redis expire error: %v\n", err)
			// Не критично, продолжаем
		}
	}

	if int(count) > uc.maxRequestsPerMinute {
		return false
	}

	return true
}
