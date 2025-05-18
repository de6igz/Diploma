package kafka_repository

import (
	"context"
	"encoding/json"
	"rule-engine-errors/internal/domain"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

// KafkaAlertDispatcher отправляет уведомления в разные Kafka-топики.
type KafkaAlertDispatcher struct {
	mailWriter     *kafka.Writer
	telegramWriter *kafka.Writer
	discordWriter  *kafka.Writer
	logger         *zerolog.Logger
}

// NewKafkaAlertDispatcher создаёт экземпляр KafkaAlertDispatcher, инициализируя kafka.Writer для каждого топика.
func NewKafkaAlertDispatcher(brokers []string, logger *zerolog.Logger) *KafkaAlertDispatcher {
	return &KafkaAlertDispatcher{
		mailWriter: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: "mail-alert-kafka-topic",
		},
		telegramWriter: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: "telegram-alert-kafka-topic",
		},
		discordWriter: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: "discord-alert-kafka-topic",
		},
		logger: logger,
	}
}

// DispatchActions отправляет каждое действие в соответствующий топик.
func (kad *KafkaAlertDispatcher) DispatchActions(ctx context.Context, e *domain.Event, acts []domain.Action) error {
	kad.logger.Info().Msgf("Dispatching %d actions", len(acts))
	for _, a := range acts {
		switch a.Type {
		case domain.ActionMail:
			kad.sendToTopic(ctx, kad.mailWriter, e, a)
		case domain.ActionTelegram:
			kad.sendToTopic(ctx, kad.telegramWriter, e, a)
		case domain.ActionDiscord:
			kad.sendToTopic(ctx, kad.discordWriter, e, a)
		default:
			kad.logger.Warn().Msgf("Unknown action type: %s", a.Type)
		}
	}
	return nil
}

func (kad *KafkaAlertDispatcher) sendToTopic(ctx context.Context, writer *kafka.Writer, e *domain.Event, a domain.Action) {
	payload := map[string]interface{}{
		"event":  e,
		"action": a,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		kad.logger.Error().Err(err).Msg("Failed to marshal payload")
		return
	}
	err = writer.WriteMessages(ctx, kafka.Message{
		Value: b,
	})
	if err != nil {
		kad.logger.Error().Err(err).Msgf("Failed to write message to topic %s", writer.Topic)
	} else {
		kad.logger.Debug().Msgf("Sent action=%s to topic=%s", a.Type, writer.Topic)
	}
}

// Close закрывает все подключения (писатели).
func (kad *KafkaAlertDispatcher) Close() error {
	var err error
	if cerr := kad.mailWriter.Close(); cerr != nil {
		err = cerr
	}
	if cerr := kad.telegramWriter.Close(); cerr != nil {
		err = cerr
	}
	if cerr := kad.discordWriter.Close(); cerr != nil {
		err = cerr
	}
	return err
}
