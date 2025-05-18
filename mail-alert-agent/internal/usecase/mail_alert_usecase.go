package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"mail-alert-agent/internal/dataproviders/email_repository"
	"mail-alert-agent/internal/dataproviders/timescale_repository"
)

// MailAlertUsecase содержит зависимости для обработки сообщения.
type MailAlertUsecase struct {
	emailRepo     email_repository.EmailRepository
	timescaleRepo timescale_repository.TimescaleRepository
	logger        *zerolog.Logger
}

// NewMailAlertUsecase создаёт экземпляр usecase.
func NewMailAlertUsecase(
	emailRepo email_repository.EmailRepository,
	timescaleRepo timescale_repository.TimescaleRepository,
	logger *zerolog.Logger,
) *MailAlertUsecase {
	return &MailAlertUsecase{
		emailRepo:     emailRepo,
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

// buildEmailMessage собирает сообщение для письма, оборачивая событие в блок кода.
func buildEmailMessage(event json.RawMessage) string {
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
	return fmt.Sprintf("%s", string(formattedBytes))
}

// ProcessMessage реализует бизнес-логику: парсинг сообщения, сборка email-сообщения,
// отправка письма через email-репозиторий с использованием context с таймаутом и логирование результата.
func (u *MailAlertUsecase) ProcessMessage(msg *sarama.ConsumerMessage) error {
	startTime := time.Now()
	u.logger.Info().Msgf("Starting processing Kafka message at offset %d", msg.Offset)

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

	if alert.Action.Type != "EMAIL" {
		u.logger.Info().Msgf("Skipping message with action type: %s", alert.Action.Type)
		return nil
	}

	emailAddress := alert.Action.Params.Value
	body := buildEmailMessage(alert.Event)
	subject := "Alert Notification"

	// Создаем context с таймаутом 15 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	u.logger.Info().Msgf("Attempting to send email to %s with subject '%s'", emailAddress, subject)

	// Отправка email в отдельной горутине с ожиданием результата или таймаута
	emailErrCh := make(chan error, 1)
	go func() {
		u.logger.Debug().Msg("Calling emailRepo.SendEmail")
		emailErrCh <- u.emailRepo.SendEmail(emailAddress, subject, body)
	}()

	select {
	case err := <-emailErrCh:
		if err != nil {
			u.logger.Error().Err(err).Msg("Failed to send email alert")
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
		u.logger.Info().Msgf("Email alert sent successfully to %s", emailAddress)
	case <-ctx.Done():
		timeoutErr := fmt.Errorf("timeout sending email alert")
		u.logger.Error().Err(timeoutErr).Msg("Email sending timed out")
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

	u.logger.Info().Msgf("Successfully processed email alert for offset %d", msg.Offset)
	logEntry := timescale_repository.LogEntry{
		KafkaTopic:     msg.Topic,
		KafkaPartition: msg.Partition,
		KafkaOffset:    msg.Offset,
		Timestamp:      time.Now(),
		Status:         "SUCCESS",
		Error:          "",
		RawMessage:     msg.Value,
	}
	_ = u.timescaleRepo.InsertLog(ctx, logEntry)
	u.logger.Info().Msgf("Processed message offset %d in %v", msg.Offset, time.Since(startTime))
	return nil
}
