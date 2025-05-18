package kafka_repository

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"telegram-alert-agent/internal/usecase"
)

// KafkaRepository инкапсулирует работу с Kafka.
type KafkaRepository struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
	logger        *zerolog.Logger
	usecase       *usecase.TelegramAlertUsecase
}

// NewKafkaRepository создаёт новый экземпляр KafkaRepository.
func NewKafkaRepository(brokers []string, consumerGroupID string, topics []string, logger *zerolog.Logger, usecase *usecase.TelegramAlertUsecase) (*KafkaRepository, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	// Отключаем авто-коммит: смещение подтверждается после успешной обработки
	cfg.Consumer.Offsets.AutoCommit.Enable = false
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Version = sarama.V2_1_0_0

	cg, err := sarama.NewConsumerGroup(brokers, consumerGroupID, cfg)
	if err != nil {
		return nil, err
	}

	return &KafkaRepository{
		consumerGroup: cg,
		topics:        topics,
		logger:        logger,
		usecase:       usecase,
	}, nil
}

// StartConsuming запускает процесс чтения сообщений из Kafka.
func (kr *KafkaRepository) StartConsuming(ctx context.Context) error {
	defer kr.consumerGroup.Close()
	handler := &consumerGroupHandler{
		usecase: kr.usecase,
		logger:  kr.logger,
	}
	for {
		if err := kr.consumerGroup.Consume(ctx, kr.topics, handler); err != nil {
			kr.logger.Error().Err(err).Msg("Error during Kafka consumption")
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type consumerGroupHandler struct {
	usecase *usecase.TelegramAlertUsecase
	logger  *zerolog.Logger
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	h.logger.Info().Msg("Kafka consumer group session setup")
	// Логируем, что агент успешно поднялся
	h.logger.Info().Msg("Telegram Alert Agent Started")
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.logger.Info().Msg("Kafka consumer group session cleanup")
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := h.usecase.ProcessMessage(msg)
		if err != nil {
			h.logger.Error().Err(err).Msgf("Error processing message at offset %d", msg.Offset)
			// Не отмечаем сообщение – оно будет переработано
			continue
		}
		session.MarkMessage(msg, "")
		h.logger.Info().Msgf("Processed message offset %d", msg.Offset)
	}
	return nil
}
