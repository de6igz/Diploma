package main

import (
	"auth-service/migrations"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/handlers"
)

//swag init -g cmd/server/main.go --output docs

// @title Auth Service API
// @version 1.0
// @description Пример Auth Service с Echo, PostgreSQL, Redis и Swagger
// @host localhost:8082
// @BasePath /
func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаемся к PostgreSQL
	if err := db.InitPostgres(cfg); err != nil {
		log.Fatalf("Failed to init postgres: %v", err)
	}
	defer db.ClosePostgres()

	// Применение миграций
	migrations.RunMigrations()

	// Подключаемся к Redis
	if err := db.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}
	defer db.CloseRedis()

	// Создаём Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Добавление настройки CORS, разрешающей все запросы и все заголовки
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
	}))

	// Регистрируем роуты аутентификации
	handlers.RegisterAuthRoutes(e, cfg)

	// Запуск сервера
	srvPort := cfg.ServicePort
	log.Printf("[INFO] Auth service running on :%s", srvPort)
	if err := e.Start(":" + srvPort); err != nil {
		log.Fatal(err)
	}
}
