package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"telegram-alert-agent/internal/dataproviders/telegram_repository"
	"telegram-alert-agent/internal/dataproviders/timescale_repository"
)

// TelegramAlertUsecase содержит зависимости для обработки сообщения.
type TelegramAlertUsecase struct {
	telegramRepo  telegram_repository.TelegramRepository
	timescaleRepo timescale_repository.TimescaleRepository
	logger        *zerolog.Logger
}

// NewTelegramAlertUsecase создаёт экземпляр usecase.
func NewTelegramAlertUsecase(
	telegramRepo telegram_repository.TelegramRepository,
	timescaleRepo timescale_repository.TimescaleRepository,
	logger *zerolog.Logger,
) *TelegramAlertUsecase {
	return &TelegramAlertUsecase{
		telegramRepo:  telegramRepo,
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

// buildTelegramMessage собирает сообщение для Telegram в формате MarkdownV2.
// Форматирует JSON с отступами и оборачивает его в блок кода.
func buildTelegramMessage(event json.RawMessage) string {
	// Форматируем JSON с отступами
	var formattedJSON map[string]interface{}
	if err := json.Unmarshal(event, &formattedJSON); err != nil {
		return "Ошибка парсинга JSON: " + err.Error()
	}

	formattedBytes, err := json.MarshalIndent(formattedJSON, "", "  ")
	if err != nil {
		return "Ошибка форматирования JSON: " + err.Error()
	}

	// Оборачиваем в блок кода с указанием языка JSON
	return fmt.Sprintf("```json\n%s\n```", string(formattedBytes))
}

// ProcessMessage реализует бизнес-логику: парсинг сообщения, сборка Telegram-сообщения,
// отправка его через репозиторий с использованием контекста с таймаутом и логирование результата.
func (u *TelegramAlertUsecase) ProcessMessage(msg *sarama.ConsumerMessage) error {
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

	if alert.Action.Type != "TELEGRAM" {
		u.logger.Info().Msgf("Skipping message with action type: %s", alert.Action.Type)
		return nil
	}

	chatID := alert.Action.Params.Value
	text := buildTelegramMessage(alert.Event)
	u.logger.Info().Msgf("Built Telegram message for chatID %s", chatID)

	// Создаем context с таймаутом 15 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	u.logger.Info().Msgf("Attempting to send Telegram message to chatID: %s", chatID)
	sendErrCh := make(chan error, 1)
	go func() {
		u.logger.Debug().Msg("Calling telegramRepo.SendMessage")
		sendErrCh <- u.telegramRepo.SendMessage(chatID, text)
	}()

	select {
	case err := <-sendErrCh:
		if err != nil {
			u.logger.Error().Err(err).Msg("Failed to send Telegram message")
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
		u.logger.Info().Msgf("Telegram message sent successfully to chatID %s", chatID)
	case <-ctx.Done():
		timeoutErr := fmt.Errorf("timeout sending Telegram message")
		u.logger.Error().Err(timeoutErr).Msg("Telegram message sending timed out")
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

	u.logger.Info().Msgf("Successfully sent Telegram message for offset %d", msg.Offset)
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
