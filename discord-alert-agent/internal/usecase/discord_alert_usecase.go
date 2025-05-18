package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"discord-alert-agent/internal/dataproviders/discord_repository"
	"discord-alert-agent/internal/dataproviders/timescale_repository"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
)

// DiscordAlertUsecase содержит зависимости для обработки сообщений.
type DiscordAlertUsecase struct {
	discordRepo   discord_repository.DiscordRepository
	timescaleRepo timescale_repository.TimescaleRepository
	logger        *zerolog.Logger
}

// NewDiscordAlertUsecase создаёт экземпляр usecase.
func NewDiscordAlertUsecase(
	discordRepo discord_repository.DiscordRepository,
	timescaleRepo timescale_repository.TimescaleRepository,
	logger *zerolog.Logger,
) *DiscordAlertUsecase {
	return &DiscordAlertUsecase{
		discordRepo:   discordRepo,
		timescaleRepo: timescaleRepo,
		logger:        logger,
	}
}

// AlertMessage описывает структуру входящего Kafka-сообщения.
type AlertMessage struct {
	Action struct {
		Type   string `json:"type"`
		Params struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"params"`
	} `json:"action"`
	Event json.RawMessage `json:"event"`
}

// buildDiscordMessage формирует сообщение для Discord, форматируя JSON с отступами
// и оборачивая его в блок кода (Discord поддерживает тройные обратные апострофы).
func buildDiscordMessage(event json.RawMessage) string {
	var formattedJSON map[string]interface{}
	if err := json.Unmarshal(event, &formattedJSON); err != nil {
		return "Ошибка парсинга JSON: " + err.Error()
	}
	formattedBytes, err := json.MarshalIndent(formattedJSON, "", "  ")
	if err != nil {
		return "Ошибка форматирования JSON: " + err.Error()
	}
	return fmt.Sprintf("```json\n%s\n```", string(formattedBytes))
}

// ProcessMessage реализует бизнес-логику: парсинг сообщения, сборка Discord-сообщения,
// отправка его через репозиторий с использованием context с таймаутом и логирование результата.
func (u *DiscordAlertUsecase) ProcessMessage(msg *sarama.ConsumerMessage) error {
	startTime := time.Now()
	u.logger.Info().Msgf("Start processing Kafka message at offset %d", msg.Offset)

	var alert AlertMessage
	err := json.Unmarshal(msg.Value, &alert)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to unmarshal Kafka message")
		logEntry := timescale_repository.LogEntry{
			KafkaTopic:     msg.Topic,
			KafkaPartition: msg.Partition,
			KafkaOffset:    msg.Offset,
			Timestamp:      time.Now(),
			Status:         "ERROR",
			Error:          fmt.Sprintf("JSON unmarshal error: %v", err),
			RawMessage:     msg.Value,
		}
		_ = u.timescaleRepo.InsertLog(context.Background(), logEntry)
		return err
	}

	if alert.Action.Type != "DISCORD" {
		u.logger.Info().Msgf("Skipping message with action type: %s", alert.Action.Type)
		return nil
	}

	channelID := alert.Action.Params.Value
	messageText := buildDiscordMessage(alert.Event)
	u.logger.Info().Msgf("Built Discord message for channelID %s", channelID)

	// Создаем context с таймаутом 15 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	u.logger.Info().Msgf("Attempting to send Discord message to channelID: %s", channelID)
	sendErrCh := make(chan error, 1)
	go func() {
		u.logger.Debug().Msg("Calling discordRepo.SendMessage")
		sendErrCh <- u.discordRepo.SendMessage(channelID, messageText)
	}()

	select {
	case err := <-sendErrCh:
		if err != nil {
			u.logger.Error().Err(err).Msg("Failed to send Discord message")
			logEntry := timescale_repository.LogEntry{
				KafkaTopic:     msg.Topic,
				KafkaPartition: msg.Partition,
				KafkaOffset:    msg.Offset,
				Timestamp:      time.Now(),
				Status:         "ERROR",
				Error:          err.Error(),
				RawMessage:     msg.Value,
			}
			_ = u.timescaleRepo.InsertLog(ctx, logEntry)
			return err
		}
		u.logger.Info().Msgf("Discord message sent successfully to channelID %s", channelID)
	case <-ctx.Done():
		timeoutErr := fmt.Errorf("timeout sending Discord message")
		u.logger.Error().Err(timeoutErr).Msg("Discord message sending timed out")
		logEntry := timescale_repository.LogEntry{
			KafkaTopic:     msg.Topic,
			KafkaPartition: msg.Partition,
			KafkaOffset:    msg.Offset,
			Timestamp:      time.Now(),
			Status:         "ERROR",
			Error:          timeoutErr.Error(),
			RawMessage:     msg.Value,
		}
		_ = u.timescaleRepo.InsertLog(ctx, logEntry)
		return timeoutErr
	}

	u.logger.Info().Msgf("Successfully sent Discord message for offset %d", msg.Offset)
	logEntry := timescale_repository.LogEntry{
		KafkaTopic:     msg.Topic,
		KafkaPartition: msg.Partition,
		KafkaOffset:    msg.Offset,
		Timestamp:      time.Now(),
		Status:         "SUCCESS",
		Error:          "",
		RawMessage:     msg.Value,
	}
	_ = u.timescaleRepo.InsertLog(context.Background(), logEntry)
	u.logger.Info().Msgf("Processed message offset %d in %v", msg.Offset, time.Since(startTime))
	return nil
}
