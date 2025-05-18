package migrations

import (
	"auth-service/internal/db"
	"github.com/pressly/goose/v3"
	"log"
)

func RunMigrations() {

	// Применяем миграции
	if err := goose.Up(db.DBConn, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("[INFO] Migrations applied successfully")
}
