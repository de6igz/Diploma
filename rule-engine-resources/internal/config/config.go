// config.go
package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Логирование
	LogLevel  string `envconfig:"LOG_LEVEL" default:"debug"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"json"` // Новое поле для формата логирования

	// MongoDB
	Mongo struct {
		User       string `envconfig:"MONGO_USER" required:"true"`
		Password   string `envconfig:"MONGO_PASSWORD" required:"true"`
		Host       string `envconfig:"MONGO_HOST" required:"true"`
		DBName     string `envconfig:"MONGO_DB" required:"true"`
		AuthSource string `envconfig:"MONGO_AUTH_SOURCE" default:"admin"`
	} `envconfig:"MONGO"`

	// Redis
	Redis struct {
		Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
		Password string `envconfig:"REDIS_PASSWORD" default:""`
		DB       int    `envconfig:"REDIS_DB" default:"0"`
	} `envconfig:"REDIS"`

	// Kafka
	Kafka struct {
		Brokers       string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
		ConsumerGroup string `envconfig:"KAFKA_CONSUMER_GROUP" default:"resource-rule-group"`
		ResourceTopic string `envconfig:"KAFKA_RESOURCE_TOPIC" default:"kafka-resource-usage-topic"`
		//ProducerTopics []string `envconfig:"KAFKA_PRODUCER_TOPICS" required:"true"` // Например: "mail,tele,disc"
	} `envconfig:"KAFKA"`

	// TimescaleDB
	Timescale struct {
		// Добавьте необходимые параметры подключения
		Host     string `envconfig:"TIMESCALE_HOST" required:"true"`
		Port     int    `envconfig:"TIMESCALE_PORT" required:"true"`
		User     string `envconfig:"TIMESCALE_USER" required:"true"`
		Password string `envconfig:"TIMESCALE_PASSWORD" required:"true"`
		DBName   string `envconfig:"TIMESCALE_DB" required:"true"`
	} `envconfig:"TIMESCALE"`

	Postgres struct {
		// Добавьте необходимые параметры подключения
		Host     string `envconfig:"POSTGRES_HOST" required:"true"`
		Port     int    `envconfig:"POSTGRES_PORT" required:"true"`
		User     string `envconfig:"POSTGRES_USER" required:"true"`
		Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
		DBName   string `envconfig:"POSTGRES_DB" required:"true"`
	} `envconfig:"POSTGRES"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}
	return &cfg, nil
}
