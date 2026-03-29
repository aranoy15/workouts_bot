# Workouts Bot

Telegram-бот для тренировок.

## Запуск

Нужны **Go** и токен бота от [@BotFather](https://t.me/BotFather).

```bash
export BOT_TOKEN="ваш_токен"
```

База — **PostgreSQL**. Параметры: `DATABASE_URL` или переменные `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE` (см. `src/config/config.go`).

Локально Postgres можно поднять из каталога `docker/` — см. `docker/DOCKER_README.md`.

```bash
go mod download
go run ./cmd/bot
```

Или сборка и запуск:

```bash
make install
make run
```

Справка по целям Makefile: `make help`.
