package main

import (
	"aletheia-public-api/internal/api/v1/events"
	"aletheia-public-api/internal/api/v1/projects"
	"aletheia-public-api/internal/api/v1/rules"
	"aletheia-public-api/internal/dataproviders/timescale"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	//"aletheia-public-api/internal/cache"
	"aletheia-public-api/internal/config"
	"aletheia-public-api/internal/dataproviders/postgres"
	//"aletheia-public-api/internal/services/middlewares"

	"aletheia-public-api/internal/transport"

	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"
)

func main() {

	// Настраиваем логгер согласно конфигурации сервиса.
	serviceCfg := config.Service()
	log.Logger = serviceCfg.Logger()
	//ctx := log.Logger.WithContext(context.Background())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	defer log.Info().Msg("goodbye")

	// Инициализация подключения к MongoDB.
	if err := postgres.InitPostgres(); err != nil {
		log.Panic().Err(err).Stack().Msg("mongo init error")
	}

	if err := timescale.Init(); err != nil {
		log.Panic().Err(err).Stack().Msg("timescale init error")
	}

	// Инициализация кэша.
	//cache.Init()

	// Создаём сервисы для историй v1 и v2.
	svcEvents := events.NewEvents()
	svcProjects := projects.NewProjects()
	svcRules := rules.NewRules()

	// Опции транспортного слоя: middlewares и регистрация сервисов.
	services := []transport.Option{
		transport.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "*",
		})),
		transport.Rules(transport.NewRules(svcRules)),
		transport.Events(transport.NewEvents(svcEvents)),
		transport.Projects(transport.NewProjects(svcProjects)),
	}

	srv := transport.New(log.Logger, services...).WithLog()

	// Запускаем endpoint для health-check.
	srv.ServeHealth(serviceCfg.HealthBind, "OK")

	// Добавляем middleware CORS для разрешения запросов с любого домена.
	// Запускаем сервер.
	go func() {

		log.Info().Str("bind", serviceCfg.Bind).Msg("listen on")
		if err := srv.Fiber().Listen(serviceCfg.Bind); err != nil {
			log.Panic().Err(err).Stack().Msg("server error")
		}
	}()

	// Если включен pprof, запускаем его.
	if serviceCfg.EnablePPROF {
		log.Info().Str("pprof", serviceCfg.BindPPROF).Msg("pprof listen on")
		go func() {
			_ = http.ListenAndServe(serviceCfg.BindPPROF, nil)
		}()
	}

	<-shutdown
}
