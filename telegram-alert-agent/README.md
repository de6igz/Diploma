# Telegram Alert Agent

Telegram Alert Agent – это Go‑сервис, который:
- Читает сообщения из Kafka-топика `telegram-alert-kafka-topic`
- Отправляет их пользователю через Telegram (с использованием библиотеки [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api))
- Логирует результаты в TimescaleDB (с использованием [zerolog](https://github.com/rs/zerolog))

## Структура проекта

. ├── cmd │ └── main.go ├── internal │ ├── config │ │ └── config.go # Конфигурация через envconfig │ ├── dataproviders │ │ ├── kafka_repository │ │ │ └── kafka_repository.go # Работа с Kafka (IBM/sarama) │ │ ├── telegram_repository │ │ │ └── telegram_repository.go # Интеграция с Telegram │ │ └── timescale_repository │ │ └── timescale_repository.go # Логирование в TimescaleDB │ └── usecase │ └── telegram_alert_usecase.go # Бизнес-логика обработки сообщений ├── go.mod └── README.md

pgsql
Копировать

## Переменные окружения

Приложение настраивается через следующие переменные окружения:

- **LOG_LEVEL** – уровень логирования (например, `debug`, `info`, `error`).
- **LOG_FORMAT** – формат логирования (`json` или `console`).

- **TELEGRAM_BOT_TOKEN** – токен для доступа к Telegram Bot API.

- **MONGO_USER**, **MONGO_PASSWORD**, **MONGO_HOST**, **MONGO_DB**, **MONGO_AUTH_SOURCE** – параметры для подключения к MongoDB.

- **REDIS_ADDR**, **REDIS_PASSWORD**, **REDIS_DB** – параметры для подключения к Redis.

- **KAFKA_BROKERS** – список Kafka-брокеров (через запятую, например, `localhost:9092`).
- **KAFKA_CONSUMER_GROUP** – идентификатор consumer group для Kafka.

- **TIMESCALE_HOST**, **TIMESCALE_PORT**, **TIMESCALE_USER**, **TIMESCALE_PASSWORD**, **TIMESCALE_DB** – параметры подключения к TimescaleDB.

## Создание таблицы в TimescaleDB

Перед запуском приложения создайте таблицу для логов:

```sql
CREATE TABLE telegram_alert_agent_logs (
    id SERIAL PRIMARY KEY,
    kafka_topic TEXT NOT NULL,
    kafka_partition INT NOT NULL,
    kafka_offset BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL,
    error TEXT,
    raw_message JSONB NOT NULL
);
Запуск проекта
Настройте переменные окружения (например, через файл .env):

env
Копировать
LOG_LEVEL=debug
LOG_FORMAT=console
TELEGRAM_BOT_TOKEN=your_telegram_bot_token

MONGO_USER=...
MONGO_PASSWORD=...
MONGO_HOST=...
MONGO_DB=...
MONGO_AUTH_SOURCE=admin

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP=telegram-alert-agent-group

TIMESCALE_HOST=localhost
TIMESCALE_PORT=5432
TIMESCALE_USER=your_timescale_user
TIMESCALE_PASSWORD=your_timescale_password
TIMESCALE_DB=your_timescale_db
Соберите и запустите приложение:

bash
Копировать
go build -o telegram-alert-agent ./cmd/main.go
./telegram-alert-agent
Зависимости
IBM/sarama – для работы с Kafka.
go-telegram-bot-api – для интеграции с Telegram.
zerolog – для логирования.
sqlx и lib/pq – для работы с TimescaleDB (PostgreSQL).
envconfig – для загрузки конфигурации из переменных окружения.
Лицензия
MIT

yaml
Копировать

---

Таким образом, весь функционал для работы с Telegram теперь находится в пакете `internal/dataproviders/telegram_repository`, а README.md обновлён и отражает новую структуру проекта.