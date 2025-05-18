package errorEvent

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type Repository interface {
	SendErrorEvent(key []byte, value []byte) (partition int32, offset int64, err error)
}

type errorRepository struct {
	producer sarama.SyncProducer
	topic    string
}

func NewErrorRepository(producer sarama.SyncProducer, topic string) Repository {
	return &errorRepository{producer: producer, topic: topic}
}

func (kr *errorRepository) SendErrorEvent(key []byte, value []byte) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: kr.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err := kr.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send error event to Kafka: %w", err)
	}
	return partition, offset, nil
}

func BuildKafkaMessage(evt ErrorEvent, userID string) ([]byte, error) {
	data := map[string]interface{}{
		"project_id":    evt.ProjectID,
		"service_name":  evt.ServiceName,
		"environment":   evt.Environment,
		"version":       evt.Version,
		"error_message": evt.ErrorMessage,
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
		"fields":        evt.Fields,
	}
	return json.Marshal(data)
}
