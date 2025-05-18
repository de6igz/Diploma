// config.go
package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"strings"
)

type Config struct {

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
		ErrorTopic    string `envconfig:"KAFKA_ERROR_TOPIC" default:"kafka-error-topic"`
		MessageTopic  string `envconfig:"KAFKA_MESSAGE_TOPIC" default:"kafka-message-topic"`
	} `envconfig:"KAFKA"`
	// Server
	Server struct {
		Port string `envconfig:"COLLECTOR_PORT" default:"8080"`
	} `envconfig:"SERVER"`

	// Rate Limiting
	RateLimiting struct {
		MaxRequestsPerMinute int `envconfig:"RATE_LIMIT_MAX_REQUESTS_PER_MINUTE" default:"60"`
	} `envconfig:"RATE_LIMITING"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}
	return &cfg, nil
}

// Дополнительная функция для получения брокеров как срез строк
func (c *Config) GetKafkaBrokers() []string {
	if c.Kafka.Brokers == "" {
		return []string{"localhost:9092"} // дефолт
	}
	return strings.Split(c.Kafka.Brokers, ",")
}
