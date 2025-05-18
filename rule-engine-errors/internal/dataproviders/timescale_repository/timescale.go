package timescale_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"rule-engine-errors/internal/config"
	"rule-engine-errors/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Подключаем драйвер pq
	"github.com/rs/zerolog"
)

// Интерфейс для репозитория
type TimescaleRepository interface {
	InsertLog(ctx context.Context, entry LogEntry) error
	Close() error
}

// Структура для логов
type LogEntry struct {
	UserID      int             `db:"user_id"`
	ServiceName string          `db:"service_name"`
	Timestamp   time.Time       `db:"timestamp"`
	Log         json.RawMessage `db:"log"`
	EventType   string          `db:"event_type"`
	UsedRules   []domain.Rule   `db:"used_rules"`
	Language    string          `db:"language"`
	ActionUsed  []domain.Action `db:"action"`
	ProjectId   string          `db:"project_id"`
	Engine      string          `db:"engine"`
}

// Реализация репозитория
type timescaleRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// Фабрика: инициализация подключения к TimescaleDB
func NewTimescaleRepository(logger *zerolog.Logger, config *config.Config) (TimescaleRepository, error) {
	user := config.Timescale.User
	pass := config.Timescale.Password
	host := config.Timescale.Host
	dbName := config.Timescale.DBName
	port := config.Timescale.Port
	if port == 0 {
		port = 5433
	}

	// Формируем DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, dbName,
	)

	// Подключаемся через sqlx
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to TimescaleDB")
		return nil, err
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		logger.Error().Err(err).Msg("Failed to ping TimescaleDB")
		return nil, err
	}

	logger.Info().Msg(fmt.Sprintf("Connected to TimescaleDB successfully %v", dsn))
	return &timescaleRepository{
		db:     db,
		logger: logger,
	}, nil
}

// Вставляет лог в таблицу logs_events
func (r *timescaleRepository) InsertLog(ctx context.Context, entry LogEntry) error {

	usedRules := r.mapUsedRules(entry.UsedRules)
	usedRulesJson, err := json.Marshal(usedRules)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to marshal used rules")
		return err
	}

	usedActionsJson, err := json.Marshal(entry.ActionUsed)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to marshal used actions")
		return err
	}

	query := `
        INSERT INTO logs_events (user_id, service_name, timestamp, log, event_type, used_rules, language,used_actions, project_id, engine) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9,$10)
    `
	_, err = r.db.ExecContext(ctx, query,
		entry.UserID,
		entry.ServiceName,
		entry.Timestamp,
		entry.Log,
		entry.EventType,
		usedRulesJson,
		entry.Language,
		usedActionsJson,
		entry.ProjectId,
		entry.Engine,
	)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to insert log into TimescaleDB")
		return err
	}
	return nil
}

// Закрывает подключение к базе
func (r *timescaleRepository) Close() error {
	return r.db.Close()
}

// UsedRule – структура для хранения id и name правила
type UsedRule struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
}

// mapUsedRules преобразует срез domain.Rule в массив объектов UsedRule
func (r *timescaleRepository) mapUsedRules(rules []domain.Rule) []UsedRule {
	usedRules := make([]UsedRule, 0, len(rules))
	for _, rule := range rules {
		usedRules = append(usedRules, UsedRule{
			RuleID:   rule.ID,
			RuleName: rule.Name,
		})
	}
	return usedRules
}
