package main

import (
	"aletheia-collector-service/internal/api/errorEvent"
	"aletheia-collector-service/internal/api/health"
	"aletheia-collector-service/internal/api/messageEvent"
	"aletheia-collector-service/internal/api/resourceEvent"
	"aletheia-collector-service/internal/config"
	"aletheia-collector-service/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	redisc "github.com/redis/go-redis/v9"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Не удалось загрузить конфигурацию: %v\n", err)
		os.Exit(1)
	}

	// Инициализация Kafka Producer
	kafkaBrokers := cfg.GetKafkaBrokers()
	kafkaProducer, err := initKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatalf("Failed to init Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Инициализируем Redis
	rdb := redisc.NewClient(&redisc.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	// Инициализируем репозитории с использованием топиков из конфигурации
	resourceKafkaRepository := resourceEvent.NewResourceRepository(kafkaProducer, cfg.Kafka.ResourceTopic, rdb)
	errorKafkaRepository := errorEvent.NewErrorRepository(kafkaProducer, cfg.Kafka.ErrorTopic)
	messageKafkaRepository := messageEvent.NewMessageRepository(kafkaProducer, cfg.Kafka.MessageTopic)

	// Инициализируем usecase
	resourceUC := resourceEvent.NewUseCase(resourceKafkaRepository, cfg.RateLimiting.MaxRequestsPerMinute)
	errorUC := errorEvent.NewUseCase(errorKafkaRepository)
	messageUC := messageEvent.NewUseCase(messageKafkaRepository)

	// Инициализируем HTTP handler
	resourceHandler := resourceEvent.NewHandler(resourceUC)
	errorHandler := errorEvent.NewHandler(errorUC)
	messageHandler := messageEvent.NewHandler(messageUC)

	// Регистрируем маршруты
	http.HandleFunc("/api/v1/resource", middleware.AuthMiddleware(resourceHandler.SendResourceEvent))
	http.HandleFunc("/api/v1/error", middleware.AuthMiddleware(errorHandler.SendErrorEvent))
	http.HandleFunc("/api/v1/event", middleware.AuthMiddleware(messageHandler.SendMessageEvent)) // todo потом поправить, чтобы в sdk и тут былло одинаково
	http.HandleFunc("/health", health.Health)

	// Запуск HTTP сервера с использованием порта из конфигурации
	port := ":" + cfg.Server.Port
	//port = ":8086"
	log.Printf("Listening on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}

func initKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	log.Printf("Kafka producer connected to brokers: %v", brokers)
	return producer, nil
}
