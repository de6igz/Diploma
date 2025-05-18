package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"aletheia-public-api/internal/config"

	_ "github.com/lib/pq"
)

var (
	// GlobalInstance — глобальное подключение к Postgres.
	GlobalInstance *sql.DB
)

// InitPostgres устанавливает соединение с Postgres, используя конфигурацию из config.Postgres().
// Строка подключения формируется следующим образом:
// "postgres://<user>:<password>@<host>:<port>/<dbName>?sslmode=disable"
func InitPostgres() error {
	cfg := config.Postgres() // Читает конфигурацию с префиксом POSTGRES (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_DB, POSTGRES_PORT)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open Postgres connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping Postgres: %w", err)
	}

	GlobalInstance = db
	log.Printf("Подключение к Postgres установлено (база: %s)", cfg.DBName)
	return nil
}
