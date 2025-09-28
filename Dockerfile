# Многоэтапный Dockerfile для workouts_bot
# Этап 1: Сборка
FROM golang:1.22-alpine AS builder

# Установка необходимых пакетов для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Создание пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bot cmd/bot/main.go

# Этап 2: Финальный образ
FROM scratch

# Копирование сертификатов и временных зон
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Копирование пользователя из builder
COPY --from=builder /etc/passwd /etc/passwd

# Копирование собранного бинарного файла
COPY --from=builder /app/bot /bot

# Переключение на непривилегированного пользователя
USER appuser

# Создание директории для логов
RUN mkdir -p /logs

# Установка рабочей директории
WORKDIR /

# Открытие порта (если потребуется для health checks)
EXPOSE 8080

# Переменные окружения по умолчанию
ENV LOG_LEVEL=info
ENV LOG_CONSOLE=true
ENV LOG_FILE_PATH=/logs/bot.log
ENV LOG_MAX_SIZE=100
ENV LOG_MAX_BACKUPS=3
ENV LOG_MAX_AGE=28
ENV LOG_COMPRESS=true

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/bot", "--health-check"] || exit 1

# Команда запуска
ENTRYPOINT ["/bot"]
