// resourceEvent/repository.go
package resourceEvent

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

// Repository – интерфейс для записи событий Resource в любое хранилище (Kafka, БД…).
type Repository interface {
	SendResourceEvent(key []byte, value []byte) (partition int32, offset int64, err error)
	GetRedisClient() *redis.Client
}

// resourceRepository – конкретная реализация для Kafka.
type resourceRepository struct {
	producer sarama.SyncProducer
	rs       *redis.Client
	topic    string
}

// NewResourceRepository – конструктор, инициализирующий репозиторий для работы с Kafka.
func NewResourceRepository(producer sarama.SyncProducer, topic string, redisClient *redis.Client) Repository {
	return &resourceRepository{
		producer: producer,
		rs:       redisClient,
		topic:    topic,
	}
}

func (kr *resourceRepository) SendResourceEvent(key []byte, value []byte) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: kr.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := kr.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send resource event to Kafka: %w", err)
	}
	return partition, offset, nil
}

// GetRedisClient возвращает клиент Redis.
func (kr *resourceRepository) GetRedisClient() *redis.Client {
	return kr.rs
}

// BuildKafkaMessage – сериализация ResourceEvent + userID (если нужно) в JSON.
func BuildKafkaMessage(evt ResourceEvent, userID string) ([]byte, error) {
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
		"fields":        evt.Fields,
		"language":      evt.Language,
	}
	return json.Marshal(data)
}
