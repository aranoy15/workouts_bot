# Makefile для workouts_bot

# Переменные
BINARY_NAME=bot
BINARY_PATH=bin/$(BINARY_NAME)
PKG_PATH=./cmd/bot

SERVICE_BINARY_NAME=service
SERVICE_BINARY_PATH=bin/$(SERVICE_BINARY_NAME)
SERVICE_PKG_PATH=./cmd/service

# Go переменные
GO=go

# Wire (совпадает с директивой go:generate в wire_gen.go)
WIRE=$(GO) run -mod=mod github.com/google/wire/cmd/wire

# Цвета для вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build build-service run run-service clean test fmt vet install dev stop logs wire wire-bot wire-service

# Помощь - показывает все доступные команды
help:
	@echo "$(BLUE)Доступные команды:$(NC)"
	@echo "  $(GREEN)build$(NC)          - Собрать бинарник бота ($(BINARY_PATH))"
	@echo "  $(GREEN)build-service$(NC)  - Собрать бинарник HTTP-сервиса ($(SERVICE_BINARY_PATH))"
	@echo "  $(GREEN)run$(NC)            - Собрать и запустить бота"
	@echo "  $(GREEN)run-service$(NC)    - Собрать и запустить сервис"
	@echo "  $(GREEN)dev$(NC)            - Запустить бота в режиме разработки"
	@echo "  $(GREEN)install$(NC)        - Установить зависимости"
	@echo "  $(GREEN)test$(NC)           - Запустить тесты"
	@echo "  $(GREEN)fmt$(NC)            - Форматировать код"
	@echo "  $(GREEN)vet$(NC)            - Проверить код на ошибки"
	@echo "  $(GREEN)clean$(NC)          - Очистить собранные файлы"
	@echo "  $(GREEN)stop$(NC)           - Остановить запущенный бот"
	@echo "  $(GREEN)logs$(NC)           - Показать логи бота"
	@echo "  $(GREEN)wire$(NC)           - Сгенерировать wire для bot и service"
	@echo "  $(GREEN)wire-bot$(NC)       - Сгенерировать wire только для cmd/bot"
	@echo "  $(GREEN)wire-service$(NC)   - Сгенерировать wire только для cmd/service"
	@echo "  $(GREEN)help$(NC)           - Показать эту справку"

# Создать директорию bin если её нет
bin:
	@mkdir -p bin

# Установить зависимости
install:
	@echo "$(BLUE)Установка зависимостей...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Зависимости установлены!$(NC)"

# Собрать бинарный файл (бот)
build: bin
	@echo "$(BLUE)Сборка $(BINARY_NAME)...$(NC)"
	$(GO) build -o $(BINARY_PATH) $(PKG_PATH)
	@echo "$(GREEN)Сборка завершена: $(BINARY_PATH)$(NC)"

# Собрать HTTP-сервис
build-service: bin
	@echo "$(BLUE)Сборка $(SERVICE_BINARY_NAME)...$(NC)"
	$(GO) build -o $(SERVICE_BINARY_PATH) $(SERVICE_PKG_PATH)
	@echo "$(GREEN)Сборка завершена: $(SERVICE_BINARY_PATH)$(NC)"

# Запустить бота
run: build
	@echo "$(BLUE)Запуск бота...$(NC)"
	@echo "$(YELLOW)Для остановки нажмите Ctrl+C$(NC)"
	./$(BINARY_PATH)

# Запустить сервис
run-service: build-service
	@echo "$(BLUE)Запуск сервиса...$(NC)"
	@echo "$(YELLOW)Для остановки нажмите Ctrl+C$(NC)"
	./$(SERVICE_BINARY_PATH)

# Сгенерировать wire (оба приложения)
wire: wire-bot wire-service
	@echo "$(GREEN)Wire: готово$(NC)"

wire-bot:
	@echo "$(BLUE)Wire: cmd/bot...$(NC)"
	$(WIRE) $(PKG_PATH)

wire-service:
	@echo "$(BLUE)Wire: cmd/service...$(NC)"
	$(WIRE) $(SERVICE_PKG_PATH)

# Запустить бота в режиме разработки (требует air)
dev: install bin
	@echo "$(BLUE)Запуск в режиме разработки...$(NC)"
	@if ! command -v air > /dev/null; then \
		echo "$(YELLOW)Установка air для автоперезагрузки...$(NC)"; \
		$(GO) install github.com/cosmtrek/air@latest; \
	fi
	@echo "$(YELLOW)Для остановки нажмите Ctrl+C$(NC)"
	air

# Запустить тесты
test:
	@echo "$(BLUE)Запуск тестов...$(NC)"
	$(GO) test -v ./...

# Форматировать код
fmt:
	@echo "$(BLUE)Форматирование кода...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)Код отформатирован!$(NC)"

# Проверить код на ошибки
vet:
	@echo "$(BLUE)Проверка кода...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)Проверка завершена!$(NC)"

# Остановить запущенный бот
stop:
	@echo "$(BLUE)Остановка бота...$(NC)"
	@pkill -f $(BINARY_NAME) || echo "$(YELLOW)Бот не запущен$(NC)"
	@echo "$(GREEN)Бот остановлен!$(NC)"

# Показать логи
logs:
	@echo "$(BLUE)Логи бота:$(NC)"
	@if [ -f logs/bot.log ]; then \
		tail -f logs/bot.log; \
	else \
		echo "$(YELLOW)Файл логов не найден$(NC)"; \
	fi

# Очистить собранные файлы
clean:
	@echo "$(BLUE)Очистка...$(NC)"
	rm -rf bin/
	@echo "$(GREEN)Очистка завершена!$(NC)"

# По умолчанию показываем помощь
.DEFAULT_GOAL := help
