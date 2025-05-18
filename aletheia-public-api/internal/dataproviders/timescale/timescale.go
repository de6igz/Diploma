package timescale

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"aletheia-public-api/internal/config"

	_ "github.com/lib/pq" // драйвер PostgreSQL, необходим для TimescaleDB
)

// GlobalInstance — глобальное подключение к TimescaleDB.
// Его можно использовать в других частях приложения для выполнения запросов к базе данных.
var GlobalInstance *sql.DB

// Init устанавливает соединение с TimescaleDB, используя конфигурацию из config.Timescale().
// Если какие-либо креды не переданы или происходит ошибка подключения, функция вернёт ошибку.
func Init() error {
	tsCfg := config.Timescale()
	// Формируем строку подключения для PostgreSQL (TimescaleDB основан на Postgres)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		tsCfg.User, tsCfg.Password, tsCfg.Host, tsCfg.Port, tsCfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open timescale connection: %w", err)
	}

	// Пингуем базу, чтобы убедиться, что соединение установлено.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping timescale: %w", err)
	}

	GlobalInstance = db
	log.Printf("Connected to TimescaleDB at %s, database: %s", tsCfg.Host, tsCfg.DBName)
	return nil
}
