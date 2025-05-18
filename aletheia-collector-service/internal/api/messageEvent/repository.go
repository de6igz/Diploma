package messageEvent

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

// Repository – общий интерфейс для отправки/сохранения MessageEvent
type Repository interface {
	SendMessageEvent(key []byte, value []byte) (partition int32, offset int64, err error)
}

// messageRepository – конкретная реализация репозитория для Kafka.
type messageRepository struct {
	producer sarama.SyncProducer
	topic    string
}

// NewMessageRepository – конструктор для работы с Kafka.
// Принимает уже инициализированный SyncProducer и название топика, куда слать сообщения этого типа.
func NewMessageRepository(producer sarama.SyncProducer, topic string) Repository {
	return &messageRepository{
		producer: producer,
		topic:    topic,
	}
}

// BuildKafkaMessage – сериализуем MessageEvent + userID (если нужно) в JSON для Kafka.
func BuildKafkaMessage(evt MessageEvent, userID string) ([]byte, error) {
	data := map[string]interface{}{
		"project_id":    evt.ProjectID,
		"service_name":  evt.ServiceName,
		"environment":   evt.Environment,
		"version":       evt.Version,
		"go_version":    evt.GoVersion,
		"os":            evt.Os,
		"arch":          evt.Arch,
		"event_type":    evt.EventType,
		"event_message": evt.EventMessage,
		"stack_trace":   evt.StackTrace,
		"tags":          evt.Tags,
		"timestamp":     evt.Timestamp,
		"context_json":  evt.ContextJson,
		"user_id":       userID,
		"language":      evt.Language,
	}
	return json.Marshal(data)
}

// SendMessageEvent – отправка сообщения в Kafka (в топик, заданный при конструировании).
func (kr *messageRepository) SendMessageEvent(key []byte, value []byte) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: kr.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err := kr.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send message event to Kafka: %w", err)
	}
	return partition, offset, nil
}
