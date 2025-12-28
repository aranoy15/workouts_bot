# Деплой в Yandex Cloud

Этот документ описывает процесс настройки автоматического деплоя бота в Yandex Cloud Serverless Containers через GitHub Actions.

## Предварительные требования

1. **Yandex Cloud аккаунт** с настроенным биллингом
2. **GitHub репозиторий** с правами администратора
3. **Yandex Cloud CLI** для первоначальной настройки

## Настройка Yandex Cloud

### 1. Создание сервисного аккаунта

```bash
# Создать сервисный аккаунт
yc iam service-account create --name workouts-bot-deployer

# Получить ID сервисного аккаунта
SA_ID=$(yc iam service-account get workouts-bot-deployer --format json | jq -r '.id')

# Назначить роли
yc resource-manager folder add-access-binding <FOLDER_ID> \
  --role container-registry.images.pusher \
  --subject serviceAccount:$SA_ID

yc resource-manager folder add-access-binding <FOLDER_ID> \
  --role serverless.containers.invoker \
  --subject serviceAccount:$SA_ID

yc resource-manager folder add-access-binding <FOLDER_ID> \
  --role serverless.containers.admin \
  --subject serviceAccount:$SA_ID

# Создать ключ для сервисного аккаунта
yc iam key create --service-account-id $SA_ID --output key.json
```

### 2. Создание Container Registry

```bash
# Создать реестр
yc container registry create --name workouts-bot-registry

# Получить ID реестра
REGISTRY_ID=$(yc container registry get workouts-bot-registry --format json | jq -r '.id')
```

### 3. Создание Serverless Container

```bash
# Создать контейнер
yc serverless container create --name workouts-bot
```

## Настройка GitHub Secrets

В настройках репозитория GitHub (Settings → Secrets and variables → Actions) добавьте следующие секреты:

### Обязательные секреты

| Секрет | Описание | Как получить |
|--------|----------|--------------|
| `YC_CLOUD_ID` | ID облака | `yc config get cloud-id` |
| `YC_FOLDER_ID` | ID папки | `yc config get folder-id` |
| `YC_REGISTRY_ID` | ID Container Registry | `yc container registry list` |
| `YC_SERVICE_ACCOUNT_KEY` | Ключ сервисного аккаунта | Содержимое `key.json` в base64: `base64 -i key.json` |
| `YC_SERVICE_ACCOUNT_ID` | ID сервисного аккаунта | `yc iam service-account get workouts-bot-deployer --format json \| jq -r '.id'` |
| `YC_CONTAINER_NAME` | Имя Serverless Container | `workouts-bot` |
| `BOT_TOKEN` | Токен Telegram бота | Получить у @BotFather |
| `DATABASE_URL` | URL базы данных | `postgres://user:password@host:port/dbname` |

### Пример получения значений

```bash
# Cloud ID
yc config get cloud-id

# Folder ID
yc config get folder-id

# Registry ID
yc container registry list --format json | jq -r '.[0].id'

# Service Account ID
yc iam service-account get workouts-bot-deployer --format json | jq -r '.id'

# Service Account Key (base64)
base64 -i key.json | tr -d '\n'
```

## Использование CI/CD

### Автоматический деплой

CI автоматически запускается при:

1. **Push в main** - деплой в production с тегом `latest`
2. **Push тега v*** - деплой в production с версионным тегом
3. **Manual trigger** - ручной запуск с выбором окружения

### Ручной запуск

1. Перейдите в GitHub Actions
2. Выберите workflow "Deploy to Yandex Cloud"
3. Нажмите "Run workflow"
4. Выберите окружение (production/staging)

### Мониторинг деплоя

Workflow включает:

- ✅ **Сборку и пуш Docker образа**
- ✅ **Деплой в Serverless Container**
- ✅ **Health check** контейнера
- ✅ **Автоматическую очистку** старых ревизий
- ✅ **Rollback** при ошибках

### Логи и отладка

```bash
# Просмотр логов контейнера
yc serverless container revision logs <REVISION_ID>

# Список ревизий
yc serverless container revision list --container-name workouts-bot

# Информация о контейнере
yc serverless container get workouts-bot
```

## Локальное тестирование

### Сборка образа

```bash
# Сборка
make docker-build

# Или напрямую
docker build -t workouts-bot .
```

### Тестирование образа

```bash
# Запуск с переменными окружения
docker run --rm \
  -e BOT_TOKEN=your_token \
  -e DATABASE_URL=your_db_url \
  -e LOG_LEVEL=debug \
  workouts-bot
```

### Пуш в registry

```bash
# Логин
make docker-login

# Пуш (требует REGISTRY_ID)
make docker-push REGISTRY_ID=your_registry_id
```

## Troubleshooting

### Частые проблемы

1. **Ошибка аутентификации**
   ```
   Error: authentication failed
   ```
   - Проверьте корректность `YC_SERVICE_ACCOUNT_KEY`
   - Убедитесь, что ключ в base64 формате

2. **Недостаточно прав**
   ```
   Error: access denied
   ```
   - Проверьте роли сервисного аккаунта
   - Убедитесь, что указан правильный `YC_FOLDER_ID`

3. **Контейнер не запускается**
   ```
   Error: container failed to start
   ```
   - Проверьте переменные окружения
   - Убедитесь, что `BOT_TOKEN` корректный
   - Проверьте доступность базы данных

4. **Health check не проходит**
   - Убедитесь, что приложение отвечает на `/health`
   - Проверьте логи контейнера
   - Увеличьте timeout в health check

### Полезные команды

```bash
# Просмотр активных ревизий
yc serverless container revision list --container-name workouts-bot

# Откат на предыдущую ревизию
yc serverless container revision deploy \
  --container-name workouts-bot \
  --revision-id <PREVIOUS_REVISION_ID>

# Просмотр логов
yc logging read --group-id <LOG_GROUP_ID> --filter 'resource.type="serverless_container"'

# Удаление старых образов
yc container image list --registry-id <REGISTRY_ID>
yc container image delete <IMAGE_ID>
```

## Безопасность

1. **Никогда не коммитьте** секреты в код
2. **Используйте GitHub Secrets** для всех чувствительных данных
3. **Регулярно ротируйте** ключи сервисных аккаунтов
4. **Ограничивайте права** сервисного аккаунта минимально необходимыми
5. **Мониторьте доступ** к ресурсам через Yandex Cloud Console

## Мониторинг и алерты

Рекомендуется настроить:

1. **Yandex Monitoring** для метрик контейнера
2. **Alerting** на ошибки деплоя
3. **Log aggregation** для централизованного логирования
4. **Health checks** для проверки доступности

## Дополнительные ресурсы

- [Yandex Cloud Serverless Containers](https://cloud.yandex.ru/docs/serverless-containers/)
- [Container Registry](https://cloud.yandex.ru/docs/container-registry/)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
