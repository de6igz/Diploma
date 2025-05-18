package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Логирование
	LogLevel  string `envconfig:"LOG_LEVEL" default:"debug"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"json"` // Возможные значения: json или console

	// Discord
	Discord struct {
		BotToken string `envconfig:"DISCORD_BOT_TOKEN" required:"true"`
	} `envconfig:"DISCORD"`

	// Kafka
	Kafka struct {
		Brokers       string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
		ConsumerGroup string `envconfig:"KAFKA_CONSUMER_GROUP" default:"resource-rule-group"`
	} `envconfig:"KAFKA"`

	// TimescaleDB
	Timescale struct {
		Host     string `envconfig:"TIMESCALE_HOST" required:"true"`
		Port     int    `envconfig:"TIMESCALE_PORT" required:"true"`
		User     string `envconfig:"TIMESCALE_USER" required:"true"`
		Password string `envconfig:"TIMESCALE_PASSWORD" required:"true"`
		DBName   string `envconfig:"TIMESCALE_DB" required:"true"`
	} `envconfig:"TIMESCALE"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}
	return &cfg, nil
}
