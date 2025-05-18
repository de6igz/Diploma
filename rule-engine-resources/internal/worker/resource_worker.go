package worker

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	// Если EvaluateRulesUseCase etc. – import
	"rule-engine-resources/internal/domain"
	"rule-engine-resources/internal/usecases"
)

type ResourceWorker struct {
	consumerGroup sarama.ConsumerGroup
	useCase       *usecases.EvaluateRulesUseCase
	logger        *zerolog.Logger
}

func NewResourceWorker(
	brokers []string,
	groupID string,
	cfg *sarama.Config,
	uc *usecases.EvaluateRulesUseCase,
	log *zerolog.Logger,
) (*ResourceWorker, error) {
	cg, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}
	return &ResourceWorker{
		consumerGroup: cg,
		useCase:       uc,
		logger:        log,
	}, nil
}

func (rw *ResourceWorker) Run(ctx context.Context, topics []string) error {
	rw.logger.Info().
		Str("groupID", rw.consumerGroup.Claims().String()).
		Msgf("ResourceWorker running on topics=%v", topics)

	for {
		if err := rw.consumerGroup.Consume(ctx, topics, rw); err != nil {
			rw.logger.Error().Err(err).Msg("ResourceWorker consume error")
			return err
		}
		if ctx.Err() != nil {
			rw.logger.Warn().Msg("ResourceWorker context canceled")
			return ctx.Err()
		}
	}
}

func (rw *ResourceWorker) Setup(_ sarama.ConsumerGroupSession) error {
	rw.logger.Debug().Msg("ResourceWorker Setup called")
	return nil
}

func (rw *ResourceWorker) Cleanup(_ sarama.ConsumerGroupSession) error {
	rw.logger.Debug().Msg("ResourceWorker Cleanup called")
	return nil
}

func (rw *ResourceWorker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		rw.logger.Debug().
			Int64("offset", msg.Offset).
			Int32("partition", claim.Partition()).
			Msg("ResourceWorker got message from resources-topic")

		// 1) Десериализуем Event (JSON)
		var evt domain.Event
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			rw.logger.Warn().Err(err).Msg("Failed to unmarshal resource Event")
			sess.MarkMessage(msg, "")
			continue
		}

		// 2) Запускаем Evaluate
		err := rw.useCase.Evaluate(sess.Context(), evt)
		if err != nil {
			rw.logger.Error().Err(err).Msg("EvaluateRulesUseCase returned error (resource event)")
		}

		// 3) Mark read
		sess.MarkMessage(msg, "")
	}
	return nil
}
