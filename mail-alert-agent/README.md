# Mail Alert Agent

Mail Alert Agent – это Go‑сервис, который:
- Читает сообщения из Kafka-топика `mail-alert-kafka-topic`
- Отправляет их пользователю по электронной почте (с использованием SMTP)
- Логирует результаты в TimescaleDB (с использованием zerolog)
- Форматирует содержимое сообщения, оборачивая его в блок кода (тройные обратные апострофы)

## Структура проекта

. ├── cmd │   └── main.go # Точка входа в приложение ├── internal │   ├── config │   │   └── config.go # Конфигурация через envconfig │   ├── dataproviders │   │   ├── kafka_repository │   │   │   └── kafka_repository.go # Работа с Kafka (IBM/sarama) │   │   ├── email_repository │   │   │   └── email_repository.go # Отправка email-алертов через SMTP │   │   └── timescale_repository │   │   └── timescale_repository.go # Логирование в TimescaleDB │   └── usecase │   └── mail_alert_usecase.go # Бизнес-логика обработки сообщений ├── go.mod └── README.md

pgsql
Копировать

## Переменные окружения

Приложение настраивается через следующие переменные окружения:

- **LOG_LEVEL** – уровень логирования (например, `debug`, `info`, `error`).
- **LOG_FORMAT** – формат логирования (`json` или `console`).

- **EMAIL_SMTP_HOST** – SMTP-сервер для отправки писем.
- **EMAIL_SMTP_PORT** – порт SMTP-сервера.
- **EMAIL_USERNAME** – имя пользователя для SMTP.
- **EMAIL_PASSWORD** – пароль для SMTP.
- **EMAIL_FROM** – email-адрес отправителя.

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

EMAIL_SMTP_HOST=smtp.example.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your_username
EMAIL_PASSWORD=your_password
EMAIL_FROM=your_email@example.com

MONGO_USER=...
MONGO_PASSWORD=...
MONGO_HOST=...
MONGO_DB=...
MONGO_AUTH_SOURCE=admin

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP=mail-alert-agent-group

TIMESCALE_HOST=localhost
TIMESCALE_PORT=5432
TIMESCALE_USER=your_timescale_user
TIMESCALE_PASSWORD=your_timescale_password
TIMESCALE_DB=your_timescale_db
Соберите и запустите приложение:

bash
Копировать
go build -o mail-alert-agent ./cmd/main.go
./mail-alert-agent
Зависимости
IBM/sarama – для работы с Kafka.
zerolog – для логирования.
sqlx и lib/pq – для работы с TimescaleDB (PostgreSQL).
envconfig – для загрузки конфигурации из переменных окружения.
Стандартная библиотека net/smtp для отправки писем.
Лицензия
MIT

yaml
Копировать

---

Таким образом, данный проект реализует сервис для отправки алертов по электронной почте. Он читает сообщения из Kafka, формирует письмо (оборачивая содержимое события в блок кода), отправляет письмо через SMTP и логирует результаты в TimescaleDB, а также логирует этапы инициализации репозиториев.