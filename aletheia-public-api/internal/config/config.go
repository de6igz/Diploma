package config

import (
	"io"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const FormatJSON = "json"

// ServiceConfig содержит параметры запуска сервиса.
type ServiceConfig struct {
	LogLevel     string `envconfig:"LOG_LEVEL" default:"debug"`
	LogFormat    string `envconfig:"LOG_FORMAT" default:"console"`
	ReportCaller bool   `envconfig:"LOG_REPORT_CALLER" default:"false"`

	Bind        string `envconfig:"BIND" default:":8085"`
	HealthBind  string `envconfig:"HEALTH_BIND" default:":9091"`
	MetricsBind string `envconfig:"METRICS_BIND" default:":9090"`

	BindPPROF   string `envconfig:"BIND_PPROF" default:":8081"`
	EnablePPROF bool   `envconfig:"ENABLE_PPROF" default:"false"`
}

var service *ServiceConfig

// Service возвращает конфигурацию сервиса с ленивой инициализацией.
func Service() ServiceConfig {
	if service != nil {
		return *service
	}
	service = &ServiceConfig{}
	if err := envconfig.Process("", service); err != nil {
		log.Fatal().Err(err).Msg("error processing Service config")
	}
	return *service
}

// Logger возвращает настроенный zerolog.Logger на основе ServiceConfig.
func (cfg ServiceConfig) Logger() zerolog.Logger {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.LogLevel); err == nil {
		level = newLevel
	}
	var out io.Writer = os.Stdout
	if cfg.LogFormat != FormatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(out).Level(level).With().Timestamp().Stack()
	if cfg.ReportCaller {
		logger = logger.Caller()
	}
	return logger.Logger()
}

// MongoConfig описывает параметры подключения к MongoDB.
// Ожидается, что все параметры будут заданы через переменные окружения с префиксом MONGO:
//
//	MONGO_USER, MONGO_PASSWORD, MONGO_HOST, MONGO_DB, MONGO_AUTH_SOURCE.
type MongoConfig struct {
	User       string `envconfig:"MONGO_USER" required:"true"`
	Password   string `envconfig:"MONGO_PASSWORD" required:"true"`
	Host       string `envconfig:"MONGO_HOST" required:"true"`
	DBName     string `envconfig:"MONGO_DB" required:"true"`
	AuthSource string `envconfig:"MONGO_AUTH_SOURCE" default:"admin"`
}

var mongoConfig *MongoConfig

// Mongo возвращает конфигурацию подключения к MongoDB.
// Для получения значений используется префикс "MONGO" (например, MONGO_USER, MONGO_PASSWORD и т.д.).
func Mongo() MongoConfig {
	if mongoConfig != nil {
		return *mongoConfig
	}
	mongoConfig = &MongoConfig{}
	if err := envconfig.Process("", mongoConfig); err != nil {
		log.Fatal().Err(err).Msg("error processing Mongo config")
	}
	return *mongoConfig
}

// TimescaleConfig описывает параметры подключения к базе данных TimescaleDB.
// Здесь переменные TIMESCALE_ADDRESS, TIMESCALE_DBNAME, TIMESCALE_USER и TIMESCALE_PASSWORD обязательны.
type TimescaleConfig struct {
	Host     string `envconfig:"TIMESCALE_HOST" required:"true"`
	Port     string `envconfig:"TIMESCALE_PORT" required:"true"`
	DBName   string `envconfig:"TIMESCALE_DB" required:"true"`
	User     string `envconfig:"TIMESCALE_USER" required:"true"`
	Password string `envconfig:"TIMESCALE_PASSWORD" required:"true"`
}

var timescaleConfig *TimescaleConfig

// Timescale возвращает конфигурацию подключения к TimescaleDB.
func Timescale() TimescaleConfig {
	if timescaleConfig != nil {
		return *timescaleConfig
	}
	timescaleConfig = &TimescaleConfig{}
	if err := envconfig.Process("", timescaleConfig); err != nil {
		log.Fatal().Err(err).Msg("error processing Timescale config")
	}
	return *timescaleConfig
}

type PostgresConfig struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	DBName   string `envconfig:"POSTGRES_DB" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
}

var postgresConfig *PostgresConfig

// Postgres возвращает конфигурацию подключения к базе данных Postgres.
func Postgres() PostgresConfig {
	if postgresConfig != nil {
		return *postgresConfig
	}
	postgresConfig = &PostgresConfig{}
	if err := envconfig.Process("", postgresConfig); err != nil {
		log.Fatal().Err(err).Msg("error processing Postgres config")
	}
	return *postgresConfig
}
