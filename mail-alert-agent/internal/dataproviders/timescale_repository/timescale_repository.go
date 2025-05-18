package timescale_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mail-alert-agent/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // драйвер PostgreSQL
	"github.com/rs/zerolog"
)

// TimescaleRepository описывает интерфейс для логирования событий.
type TimescaleRepository interface {
	InsertLog(ctx context.Context, entry LogEntry) error
	Close() error
}

// LogEntry описывает структуру записи в таблице логов.
type LogEntry struct {
	KafkaTopic     string          `db:"kafka_topic"`
	KafkaPartition int32           `db:"kafka_partition"`
	KafkaOffset    int64           `db:"kafka_offset"`
	Timestamp      time.Time       `db:"timestamp"`
	Status         string          `db:"status"` // "SUCCESS" или "ERROR"
	Error          string          `db:"error"`  // текст ошибки, если есть
	RawMessage     json.RawMessage `db:"raw_message"`
}

type timescaleRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// NewTimescaleRepository инициализирует подключение к TimescaleDB.
func NewTimescaleRepository(logger *zerolog.Logger, cfg *config.Config) (TimescaleRepository, error) {
	user := cfg.Timescale.User
	pass := cfg.Timescale.Password
	host := cfg.Timescale.Host
	dbName := cfg.Timescale.DBName
	port := cfg.Timescale.Port
	if port == 0 {
		port = 5433
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to TimescaleDB")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error().Err(err).Msg("Failed to ping TimescaleDB")
		return nil, err
	}

	logger.Info().Msgf("Connected to TimescaleDB successfully: %s", dsn)
	return &timescaleRepository{
		db:     db,
		logger: logger,
	}, nil
}

// InsertLog вставляет запись в таблицу логов.
func (r *timescaleRepository) InsertLog(ctx context.Context, entry LogEntry) error {
	query := `
        INSERT INTO mail_alert_agent_logs (kafka_topic, kafka_partition, kafka_offset, timestamp, status, error, raw_message)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(ctx, query,
		entry.KafkaTopic,
		entry.KafkaPartition,
		entry.KafkaOffset,
		entry.Timestamp,
		entry.Status,
		entry.Error,
		entry.RawMessage,
	)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to insert log into TimescaleDB")
		return err
	}
	return nil
}

// Close закрывает подключение к базе.
func (r *timescaleRepository) Close() error {
	return r.db.Close()
}
