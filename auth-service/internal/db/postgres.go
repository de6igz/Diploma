package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"auth-service/internal/config"
)

var DBConn *sql.DB

func InitPostgres(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPassword, cfg.PGDB)

	var err error
	DBConn, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	return DBConn.Ping()
}

func ClosePostgres() {
	if DBConn != nil {
		DBConn.Close()
	}
}
