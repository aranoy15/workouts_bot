# Workouts Bot

Telegram бот для управления тренировками и отслеживания прогресса.

## 🚀 Быстрый старт

### С SQLite (для разработки)
```bash
# Установить зависимости
go mod tidy

# Запустить бота
go run cmd/bot/main.go
```

### С PostgreSQL (для продакшена)
```bash
# Перейти в папку Docker
cd docker

# Запустить PostgreSQL
docker-compose up -d

# Вернуться в корень и запустить бота
cd ..
go run cmd/bot/main.go
```

## 📁 Структура проекта

```
workouts_bot/
├── cmd/bot/           # Точка входа приложения
├── src/
│   ├── bot/           # Логика бота
│   ├── config/        # Конфигурация
│   └── database/      # Модели и подключение к БД
├── pkg/logger/        # Логирование
├── docker/           # Docker конфигурация
│   ├── docker-compose.yml
│   ├── init.sql
│   └── DOCKER_README.md
└── .env              # Переменные окружения
```

## 🔧 Конфигурация

Скопируйте `.env.example` в `.env` и настройте переменные:

```bash
# Bot configuration
BOT_TOKEN=your_bot_token_here

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=workouts_user
DB_PASSWORD=workouts_password
DB_NAME=workouts_db
DB_SSL_MODE=disable
```

## 📊 База данных

- **SQLite**: Для разработки (файл `workouts.db`)
- **PostgreSQL**: Для продакшена (см. `docker/` папку)

## 🐳 Docker

Подробные инструкции по работе с Docker находятся в папке `docker/`:
- `docker/DOCKER_README.md` - полное руководство
- `docker/docker-compose.yml` - конфигурация PostgreSQL
- `docker/init.sql` - инициализация базы данных
