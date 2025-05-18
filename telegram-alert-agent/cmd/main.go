package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"telegram-alert-agent/internal/config"
	"telegram-alert-agent/internal/dataproviders/kafka_repository"
	"telegram-alert-agent/internal/dataproviders/telegram_repository"
	"telegram-alert-agent/internal/dataproviders/timescale_repository"
	"telegram-alert-agent/internal/usecase"

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
	switch strings.ToLower(cfg.LogFormat) {
	case "human_read":
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(level).
			With().
			Timestamp().
			Logger()
	case "json":
		fallthrough
	default:
		log.Logger = zerolog.New(os.Stderr).
			Level(level).
			With().
			Timestamp().
			Logger()
	}

	// Инициализируем TimescaleDB репозиторий
	tsRepo, err := timescale_repository.NewTimescaleRepository(&log.Logger, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Timescale repository")
	}
	defer tsRepo.Close()
	log.Info().Msg("Timescale repository initialized")

	// Инициализируем Telegram репозиторий
	telegramRepo, err := telegram_repository.NewTelegramRepository(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Telegram repository")
	}
	log.Info().Msg("Telegram repository initialized")
	// Запускаем прослушивание входящих команд бота (например, /get_chat_id)
	go telegramRepo.StartCommandListener()

	// Создаём usecase для обработки Telegram alert
	alertUsecase := usecase.NewTelegramAlertUsecase(telegramRepo, tsRepo, &log.Logger)

	// Инициализируем Kafka репозиторий
	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	// Для telegram alert используем топик "telegram-alert-kafka-topic"
	topics := []string{"telegram-alert-kafka-topic"}
	kafkaRepo, err := kafka_repository.NewKafkaRepository(brokers, cfg.Kafka.ConsumerGroup, topics, &log.Logger, alertUsecase)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka repository")
	}
	log.Info().Msg("Kafka repository initialized")

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

	log.Info().Msg("Telegram alert agent is shutting down")
}
