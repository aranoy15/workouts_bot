# Makefile для workouts_bot

# Переменные
BINARY_NAME=bot
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/bot/main.go

# Go переменные
GO=go

# Цвета для вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run clean test fmt vet install dev stop logs docker-build docker-push docker-login docker-registry-list docker-registry-create docker-image-list

# Помощь - показывает все доступные команды
help:
	@echo "$(BLUE)Доступные команды:$(NC)"
	@echo "  $(GREEN)build$(NC)     - Собрать бинарный файл"
	@echo "  $(GREEN)run$(NC)       - Запустить бота"
	@echo "  $(GREEN)dev$(NC)       - Запустить бота в режиме разработки"
	@echo "  $(GREEN)install$(NC)   - Установить зависимости"
	@echo "  $(GREEN)test$(NC)      - Запустить тесты"
	@echo "  $(GREEN)fmt$(NC)       - Форматировать код"
	@echo "  $(GREEN)vet$(NC)       - Проверить код на ошибки"
	@echo "  $(GREEN)clean$(NC)     - Очистить собранные файлы"
	@echo "  $(GREEN)stop$(NC)      - Остановить запущенный бот"
	@echo "  $(GREEN)logs$(NC)      - Показать логи бота"
	@echo "  $(GREEN)docker-login$(NC) - Войти в Yandex Cloud Container Registry"
	@echo "  $(GREEN)docker-build$(NC) - Собрать Docker образ"
	@echo "  $(GREEN)docker-push$(NC) - Загрузить Docker образ в registry"
	@echo "  $(GREEN)docker-registry-list$(NC) - Показать список реестров"
	@echo "  $(GREEN)docker-registry-create$(NC) - Создать новый реестр"
	@echo "  $(GREEN)docker-image-list$(NC) - Показать список образов в registry"
	@echo "  $(GREEN)help$(NC)      - Показать эту справку"

# Создать директорию bin если её нет
bin:
	@mkdir -p bin

# Установить зависимости
install:
	@echo "$(BLUE)Установка зависимостей...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Зависимости установлены!$(NC)"

# Собрать бинарный файл
build: bin
	@echo "$(BLUE)Сборка $(BINARY_NAME)...$(NC)"
	$(GO) build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Сборка завершена: $(BINARY_PATH)$(NC)"

# Запустить бота
run: build
	@echo "$(BLUE)Запуск бота...$(NC)"
	@echo "$(YELLOW)Для остановки нажмите Ctrl+C$(NC)"
	./$(BINARY_PATH)

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

# Войти в Yandex Cloud Container Registry
docker-login:
	@echo "$(BLUE)Вход в Yandex Cloud Container Registry...$(NC)"
	@echo "$(YELLOW)Убедитесь, что вы выполнили 'yc init' ранее$(NC)"
	docker login cr.yandex --username iam --password $$(yc iam create-token)
	@echo "$(GREEN)Вход выполнен успешно!$(NC)"

# Собрать Docker образ
docker-build:
	@echo "$(BLUE)Сборка Docker образа...$(NC)"
	docker build -t workouts-bot .
	@echo "$(GREEN)Docker образ собран!$(NC)"

# Загрузить Docker образ в registry (требует переменную REGISTRY_ID)
docker-push: docker-build
	@echo "$(BLUE)Загрузка образа в Yandex Cloud Registry...$(NC)"
	@if [ -z "$(REGISTRY_ID)" ]; then \
		echo "$(YELLOW)REGISTRY_ID не указан, пытаюсь получить автоматически...$(NC)"; \
		REGISTRY_ID=$$(yc container registry list --format=json | jq -r '.[0].id' 2>/dev/null); \
		if [ -z "$$REGISTRY_ID" ] || [ "$$REGISTRY_ID" = "null" ]; then \
			echo "$(RED)Ошибка: Не удалось получить ID реестра автоматически$(NC)"; \
			echo "$(YELLOW)Пожалуйста, укажите ID реестра через переменную REGISTRY_ID$(NC)"; \
			echo "$(YELLOW)Или создайте новый реестр: make docker-registry-create$(NC)"; \
			echo "$(YELLOW)Пример: make docker-push REGISTRY_ID=your-registry-id$(NC)"; \
			exit 1; \
		else \
			echo "$(GREEN)Найден реестр с ID: $$REGISTRY_ID$(NC)"; \
			echo "$(BLUE)Тегирование образа...$(NC)"; \
			docker tag workouts-bot cr.yandex/$$REGISTRY_ID/workouts-bot:latest; \
			echo "$(BLUE)Загрузка образа в registry...$(NC)"; \
			docker push cr.yandex/$$REGISTRY_ID/workouts-bot:latest; \
		fi; \
	else \
		echo "$(GREEN)Используется указанный ID реестра: $(REGISTRY_ID)$(NC)"; \
		echo "$(BLUE)Тегирование образа...$(NC)"; \
		docker tag workouts-bot cr.yandex/$(REGISTRY_ID)/workouts-bot:latest; \
		echo "$(BLUE)Загрузка образа в registry...$(NC)"; \
		docker push cr.yandex/$(REGISTRY_ID)/workouts-bot:latest; \
	fi
	@echo "$(GREEN)Образ успешно загружен в registry!$(NC)"

# Показать список реестров
docker-registry-list:
	@echo "$(BLUE)Список реестров в Yandex Cloud:$(NC)"
	@yc container registry list || echo "$(RED)Ошибка: не удалось получить список реестров. Убедитесь, что вы вошли в Yandex Cloud (yc init)$(NC)"

# Создать новый реестр
docker-registry-create:
	@echo "$(BLUE)Создание нового реестра...$(NC)"
	@yc container registry create --name workouts-bot-registry || echo "$(RED)Ошибка: не удалось создать реестр$(NC)"
	@echo "$(GREEN)Реестр создан! Используйте 'make docker-registry-list' для получения ID$(NC)"

# Показать список образов в registry
docker-image-list:
	@echo "$(BLUE)Список образов в Yandex Cloud Registry:$(NC)"
	@yc container image list || echo "$(RED)Ошибка: не удалось получить список образов. Убедитесь, что вы вошли в Yandex Cloud (yc init)$(NC)"

# По умолчанию показываем помощь
.DEFAULT_GOAL := help
