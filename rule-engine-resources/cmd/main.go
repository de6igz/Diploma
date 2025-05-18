package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"rule-engine-resources/internal"
	"rule-engine-resources/internal/config"
	kafkaRepository "rule-engine-resources/internal/dataproviders/kafka_repository"
	postgres "rule-engine-resources/internal/dataproviders/postgres_repository"
	redisRepository "rule-engine-resources/internal/dataproviders/redis_repository"
	timescaleRepository "rule-engine-resources/internal/dataproviders/timescale_repository"
	"rule-engine-resources/internal/usecases"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Не удалось загрузить конфигурацию: %v\n", err)
		os.Exit(1)
	}

	// Настройка логирования
	var level zerolog.Level
	switch cfg.LogLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}

	var logger zerolog.Logger
	switch strings.ToLower(cfg.LogFormat) {
	case "human_read":
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(level).
			With().
			Timestamp().
			Logger()
	default:
		logger = zerolog.New(os.Stderr).
			Level(level).
			With().
			Timestamp().
			Logger()
	}

	logger.Info().Msgf("Starting Resource Rule Engine Resources with log-level=%s and log-format=%s", cfg.LogLevel, cfg.LogFormat)

	//// Подключение к MongoDB (с аутентификацией)
	//client, db, err := initMongoWithAuth(cfg, logger)
	//if err != nil {
	//	logger.Fatal().Err(err).Msg("Failed to init Mongo")
	//}
	//defer client.Disconnect(context.Background())
	//
	//// Инициализируем репозиторий правил из коллекции "rules_resources"
	//ruleRepo := mongoRepository.NewMongoRuleRepository(db, "rules_resources", logger)
	ruleRepo, err := postgres.NewPostgresRuleRepository(&logger, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to init Postgres repo")
	}
	// Подключаемся к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()
	repeatCounter := redisRepository.NewRedisRepeatCounter(rdb, 300, &logger)
	redisCache := redisRepository.NewRedisRuleCache(rdb, internal.RulesCacheTTL, &logger)

	// Разбиваем список брокеров (ожидается, что они разделены запятыми)
	kafkaBrokers := strings.Split(cfg.Kafka.Brokers, ",")

	// Инициализация Kafka Alert Dispatcher (отправка уведомлений в топики MAIL/TELEGRAM/DISCORD)
	dispatcher := kafkaRepository.NewKafkaAlertDispatcher(kafkaBrokers, &logger)
	defer dispatcher.Close()

	// Инициализация TimescaleDB репозитория
	timeScaleRepo, err := timescaleRepository.NewTimescaleRepository(&logger, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to init TimescaleDB repo")
	}
	defer timeScaleRepo.Close()

	// Собираем use case для оценки правил
	evalUC := usecases.NewEvaluateRulesUseCase(
		ruleRepo,
		timeScaleRepo,
		dispatcher,
		repeatCounter,
		redisCache,
		&logger,
	)

	// Инициализируем Kafka consumer с ручным коммитом (CommitInterval: 0)
	consumer := kafkaRepository.NewRuleEngineConsumer(
		kafkaBrokers,
		cfg.Kafka.ConsumerGroup,
		cfg.Kafka.ResourceTopic,
		evalUC,
		&logger,
	)
	defer consumer.Close()

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем consumer в отдельной горутине
	go func() {
		if err := consumer.Run(ctx); err != nil {
			logger.Error().Err(err).Msg("Resource consumer stopped with error")
		}
	}()

	logger.Info().Msg("Resource Rule Engine running. Waiting for signal...")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info().Msg("Shutting down Resource Rule Engine gracefully...")
}

// initMongoWithAuth устанавливает соединение с MongoDB с аутентификацией.
func initMongoWithAuth(cfg *config.Config, logger zerolog.Logger) (*mongo.Client, *mongo.Database, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=%s",
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.DBName,
		cfg.Mongo.AuthSource,
	)
	logger.Debug().Msgf("Mongo URI: %s", mongoURI)

	opts := options.Client().ApplyURI(mongoURI)
	opts.SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, nil, err
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, nil, err
	}
	logger.Info().Msg("Connected to MongoDB successfully")

	return client, client.Database(cfg.Mongo.DBName), nil
}
