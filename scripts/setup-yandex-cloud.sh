#!/bin/bash

# Скрипт для настройки Yandex Cloud ресурсов для деплоя workouts-bot
# Требует установленный и настроенный yc CLI

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для вывода сообщений
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка установки yc CLI
check_yc_cli() {
    if ! command -v yc &> /dev/null; then
        log_error "Yandex Cloud CLI не установлен"
        log_info "Установите его: curl -sSL https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash"
        exit 1
    fi

    log_success "Yandex Cloud CLI найден"
}

# Проверка аутентификации
check_auth() {
    if ! yc config get token &> /dev/null && ! yc config get service-account-key &> /dev/null; then
        log_error "Не выполнена аутентификация в Yandex Cloud"
        log_info "Выполните: yc init"
        exit 1
    fi

    log_success "Аутентификация в Yandex Cloud выполнена"
}

# Получение текущих настроек
get_current_config() {
    CLOUD_ID=$(yc config get cloud-id 2>/dev/null || echo "")
    FOLDER_ID=$(yc config get folder-id 2>/dev/null || echo "")

    if [[ -z "$CLOUD_ID" || -z "$FOLDER_ID" ]]; then
        log_error "Cloud ID или Folder ID не настроены"
        log_info "Выполните: yc init"
        exit 1
    fi

    log_info "Cloud ID: $CLOUD_ID"
    log_info "Folder ID: $FOLDER_ID"
}

# Создание сервисного аккаунта
create_service_account() {
    local sa_name="workouts-bot-deployer"

    log_info "Создание сервисного аккаунта: $sa_name"

    # Проверяем, существует ли уже сервисный аккаунт
    if yc iam service-account get "$sa_name" &> /dev/null; then
        log_warning "Сервисный аккаунт $sa_name уже существует"
        SA_ID=$(yc iam service-account get "$sa_name" --format json | jq -r '.id')
    else
        yc iam service-account create --name "$sa_name" --description "Service account for workouts-bot deployment"
        SA_ID=$(yc iam service-account get "$sa_name" --format json | jq -r '.id')
        log_success "Сервисный аккаунт создан: $SA_ID"
    fi

    # Назначение ролей
    log_info "Назначение ролей сервисному аккаунту..."

    local roles=(
        "container-registry.images.pusher"
        "serverless.containers.invoker"
        "serverless.containers.admin"
        "iam.serviceAccounts.user"
    )

    for role in "${roles[@]}"; do
        yc resource-manager folder add-access-binding "$FOLDER_ID" \
            --role "$role" \
            --subject "serviceAccount:$SA_ID" || log_warning "Роль $role уже назначена"
    done

    log_success "Роли назначены сервисному аккаунту"
}

# Создание ключа для сервисного аккаунта
create_service_account_key() {
    log_info "Создание ключа для сервисного аккаунта..."

    local key_file="workouts-bot-sa-key.json"

    if [[ -f "$key_file" ]]; then
        log_warning "Файл ключа $key_file уже существует"
        read -p "Перезаписать? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Пропускаем создание ключа"
            return
        fi
    fi

    yc iam key create --service-account-id "$SA_ID" --output "$key_file"

    # Кодируем ключ в base64 для GitHub Secrets
    SA_KEY_BASE64=$(base64 -i "$key_file" | tr -d '\n')

    log_success "Ключ создан: $key_file"
    log_info "Base64 ключ для GitHub Secrets сохранен в переменную SA_KEY_BASE64"
}

# Создание Container Registry
create_registry() {
    local registry_name="workouts-bot-registry"

    log_info "Создание Container Registry: $registry_name"

    # Проверяем, существует ли уже реестр
    if yc container registry get "$registry_name" &> /dev/null; then
        log_warning "Реестр $registry_name уже существует"
        REGISTRY_ID=$(yc container registry get "$registry_name" --format json | jq -r '.id')
    else
        yc container registry create --name "$registry_name"
        REGISTRY_ID=$(yc container registry get "$registry_name" --format json | jq -r '.id')
        log_success "Реестр создан: $REGISTRY_ID"
    fi
}

# Создание Serverless Container
create_serverless_container() {
    local container_name="workouts-bot"

    log_info "Создание Serverless Container: $container_name"

    # Проверяем, существует ли уже контейнер
    if yc serverless container get "$container_name" &> /dev/null; then
        log_warning "Контейнер $container_name уже существует"
    else
        yc serverless container create --name "$container_name"
        log_success "Serverless Container создан: $container_name"
    fi

    CONTAINER_NAME="$container_name"
}

# Вывод информации для GitHub Secrets
print_github_secrets() {
    log_info "Настройте следующие секреты в GitHub Actions:"
    echo
    echo "YC_CLOUD_ID=$CLOUD_ID"
    echo "YC_FOLDER_ID=$FOLDER_ID"
    echo "YC_REGISTRY_ID=$REGISTRY_ID"
    echo "YC_SERVICE_ACCOUNT_ID=$SA_ID"
    echo "YC_SERVICE_ACCOUNT_KEY=$SA_KEY_BASE64"
    echo "YC_CONTAINER_NAME=$CONTAINER_NAME"
    echo
    log_warning "Также не забудьте добавить:"
    echo "BOT_TOKEN=<your_telegram_bot_token>"
    echo "DATABASE_URL=<your_database_url>"
    echo
}

# Создание .env файла для локальной разработки
create_env_file() {
    local env_file=".env.example"

    log_info "Создание примера .env файла: $env_file"

    cat > "$env_file" << EOF
# Telegram Bot Configuration
BOT_TOKEN=your_telegram_bot_token_here

# Database Configuration
DATABASE_URL=postgres://user:password@localhost:5432/workouts_db?sslmode=disable

# Alternative database configuration (if DATABASE_URL is not used)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=workouts_db
DB_SSL_MODE=disable

# Webhook Configuration (for production)
WEBHOOK_ENABLED=true
WEBHOOK_URL=https://your-serverless-container-url.yandexcloud.net
WEBHOOK_PATH=/webhook
PORT=8080

# Logging Configuration
LOG_LEVEL=info
LOG_CONSOLE=true
LOG_FILE_PATH=/logs/bot.log
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=3
LOG_MAX_AGE=28
LOG_COMPRESS=true
LOG_JSON_FORMAT=false

# Yandex Cloud Configuration (for local deployment)
YC_CLOUD_ID=$CLOUD_ID
YC_FOLDER_ID=$FOLDER_ID
YC_REGISTRY_ID=$REGISTRY_ID
YC_SERVICE_ACCOUNT_ID=$SA_ID
YC_CONTAINER_NAME=$CONTAINER_NAME
EOF

    log_success "Пример .env файла создан: $env_file"
}

# Основная функция
main() {
    log_info "Настройка Yandex Cloud для workouts-bot"
    echo

    check_yc_cli
    check_auth
    get_current_config

    echo
    create_service_account
    create_service_account_key
    create_registry
    create_serverless_container

    echo
    log_success "Настройка Yandex Cloud завершена!"
    echo

    print_github_secrets
    create_env_file

    echo
    log_info "Следующие шаги:"
    echo "1. Добавьте секреты в GitHub Actions (Settings → Secrets and variables → Actions)"
    echo "2. Скопируйте .env.example в .env и заполните необходимые значения"
    echo "3. Запустите деплой через GitHub Actions или локально"
    echo
    log_success "Готово! 🚀"
}

# Запуск скрипта
main "$@"
