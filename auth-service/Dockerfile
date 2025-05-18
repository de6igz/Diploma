# Используем официальный образ Go для сборки
FROM golang:1.23.4 AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем всё остальное
COPY . .

# Собираем бинарник
RUN go build -o auth_service ./cmd/server

# Минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/auth_service .
COPY --from=builder /app/migrations/ /root/migrations/

# Открываем порт 8082 (по умолчанию)
EXPOSE 8082

CMD ["./auth_service"]
