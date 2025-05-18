package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"mail-alert-agent/internal/config"
	"mail-alert-agent/internal/dataproviders/email_repository"
	"mail-alert-agent/internal/dataproviders/kafka_repository"
	"mail-alert-agent/internal/dataproviders/timescale_repository"
	"mail-alert-agent/internal/usecase"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Настраиваем zerolog
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	if cfg.LogFormat == "human_read" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Инициализируем TimescaleDB репозиторий
	tsRepo, err := timescale_repository.NewTimescaleRepository(&log.Logger, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Timescale repository")
	}
	defer tsRepo.Close()
	log.Info().Msg("Timescale repository initialized")

	// Инициализируем Email репозиторий
	emailRepo, err := email_repository.NewEmailRepository(
		&log.Logger,
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.From,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Email repository")
	}
	log.Info().Msg("Email repository initialized")

	// Создаём usecase для обработки email alert
	alertUsecase := usecase.NewMailAlertUsecase(emailRepo, tsRepo, &log.Logger)

	// Инициализируем Kafka репозиторий
	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	// Для email alert используем топик "mail-alert-kafka-topic"
	topics := []string{"mail-alert-kafka-topic"}
	kafkaRepo, err := kafka_repository.NewKafkaRepository(brokers, cfg.Kafka.ConsumerGroup, topics, &log.Logger, alertUsecase)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka repository")
	}
	log.Info().Msg("Kafka repository initialized")

	// Логируем, что агент успешно поднялся
	log.Info().Msg("Mail Alert Agent Started")

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		cancel()
	}()

	// Запускаем чтение сообщений из Kafka
	err = kafkaRepo.StartConsuming(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error during Kafka consumption")
	}

	log.Info().Msg("Mail alert agent is shutting down")
}
