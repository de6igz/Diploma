package messageEvent

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
)

// UseCase – интерфейс бизнес-логики для MessageEvent.
type UseCase interface {
	ProcessMessageEvent(evt MessageEvent, userID string) error
}

// useCaseImpl – реализация UseCase. Хранит ссылку на Repository.
type useCaseImpl struct {
	repo Repository
}

// NewUseCase – конструктор.
func NewUseCase(repo Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

// ProcessMessageEvent – основная бизнес-логика обработки MessageEvent.
// Генерирует ключ для Kafka, формирует JSON и отправляет в репозиторий.
func (uc *useCaseImpl) ProcessMessageEvent(evt MessageEvent, userID string) error {
	// 1. Генерируем ключ (можно использовать ID сообщения, UUID и т.д.)
	key, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed to generate key for message event: %v", err)
	}

	// 2. Строим JSON для Kafka
	value, err := BuildKafkaMessage(evt, userID)
	if err != nil {
		return fmt.Errorf("failed to build message for message event: %v", err)
	}

	// 3. Отправляем в репозиторий (Kafka)
	partition, offset, err := uc.repo.SendMessageEvent([]byte(key), value)
	if err != nil {
		return err
	}

	// 4. (Опционально) логируем результат
	fmt.Printf("[MessageEvent] Sent to topic '%s' partition=%d, offset=%d\n", "kafka-message-topic", partition, offset)
	return nil
}
