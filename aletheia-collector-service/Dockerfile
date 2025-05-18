# ====== STAGE 1: Сборка Go-приложения ======
FROM golang:1.23.4 AS builder

WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Собираем бинарник (статическая линковка, если нужна)
RUN CGO_ENABLED=0 go build -o /app/collector ./cmd/main.go

# ====== STAGE 2: Финальный контейнер ======
FROM alpine:latest

WORKDIR /app

# Устанавливаем библиотеки, если нужны
RUN apk --no-cache add ca-certificates libc6-compat

# Копируем бинарник
COPY --from=builder /app/collector /app/collector

# Даем права на выполнение
RUN chmod +x /app/collector


# Запускаем приложение
CMD ["/app/collector"]
