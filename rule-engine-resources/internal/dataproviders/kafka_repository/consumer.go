package kafka_repository

import (
	"context"
	"encoding/json"
	"rule-engine-resources/internal/domain"
	"rule-engine-resources/internal/usecases"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

// RuleEngineConsumer читает сообщения из Kafka и передаёт события в EvaluateRulesUseCase.
type RuleEngineConsumer struct {
	reader  *kafka.Reader
	useCase *usecases.EvaluateRulesUseCase
	logger  *zerolog.Logger
}

// NewRuleEngineConsumer создаёт нового потребителя с ручным коммитом (CommitInterval: 0).
func NewRuleEngineConsumer(brokers []string, groupID, topic string, uc *usecases.EvaluateRulesUseCase, logger *zerolog.Logger) *RuleEngineConsumer {
	readerConfig := kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 0,    // ручной коммит
	}
	reader := kafka.NewReader(readerConfig)
	return &RuleEngineConsumer{
		reader:  reader,
		useCase: uc,
		logger:  logger,
	}
}

// Run запускает бесконечный цикл чтения сообщений.
func (rec *RuleEngineConsumer) Run(ctx context.Context) error {
	rec.logger.Info().Msg("Started resources worker")
	for {
		m, err := rec.reader.ReadMessage(ctx)
		if err != nil {
			rec.logger.Error().Err(err).Msg("Error reading Kafka message")
			return err
		}
		var evt domain.Event
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			rec.logger.Warn().Err(err).Msg("Failed to unmarshal event from Kafka")
			continue
		}
		rec.logger.Debug().Msgf("Received Event: service=%s environment=%s level=%s", evt.ServiceName, evt.Environment, evt.EventType)

		// Обработка события и отправка сообщений в соответствующие топики
		if err := rec.useCase.Evaluate(ctx, &evt); err != nil {
			rec.logger.Error().Err(err).Msg("EvaluateRulesUseCase returned error")
			continue
		}

		// После успешной отправки сообщений выполняем ручной коммит
		if err := rec.reader.CommitMessages(ctx, m); err != nil {
			rec.logger.Error().Err(err).Msg("Failed to commit message")
		} else {
			rec.logger.Info().Msg("Message committed successfully")
		}
	}
}

// Close закрывает потребителя.
func (rec *RuleEngineConsumer) Close() error {
	return rec.reader.Close()
}
