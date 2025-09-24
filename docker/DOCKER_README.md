# Workouts Bot - Docker Compose для PostgreSQL

## 🚀 Быстрый старт

### 1. Запуск PostgreSQL
```bash
# Перейти в папку docker
cd docker

# Запустить PostgreSQL
docker-compose up -d

# Проверить статус
docker-compose ps
```

### 2. Подключение к базе данных
- **PostgreSQL**: `localhost:5432`

### 3. Запуск бота
```bash
# Вернуться в корень проекта
cd ..

# Установить зависимости
go mod tidy

# Запустить бота
go run cmd/bot/main.go
```

## 📊 Управление базой данных

### Подключение через psql
```bash
# Подключиться к PostgreSQL
docker exec -it workouts_postgres psql -U workouts_user -d workouts_db
```

### Подключение через другие клиенты
Вы можете использовать любой PostgreSQL клиент:
- **psql** (командная строка)
- **DBeaver** (GUI)
- **DataGrip** (JetBrains)
- **pgAdmin** (если установлен отдельно)

Настройки подключения:
- **Host**: `localhost`
- **Port**: `5432`
- **Database**: `workouts_db`
- **Username**: `workouts_user`
- **Password**: `workouts_password`

## 🔧 Конфигурация

### Переменные окружения (.env)
```bash
# Database configuration for PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=workouts_user
DB_PASSWORD=workouts_password
DB_NAME=workouts_db
DB_SSL_MODE=disable
```

### Полезные команды
```bash
# Остановить сервисы
docker-compose down

# Остановить и удалить данные
docker-compose down -v

# Просмотр логов
docker-compose logs postgres

# Перезапуск сервисов
docker-compose restart
```

## 📁 Структура файлов
- `docker-compose.yml` - конфигурация Docker Compose
- `init.sql` - скрипт инициализации базы данных
- `.env` - переменные окружения
