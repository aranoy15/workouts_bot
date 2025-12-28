# Быстрый старт - Деплой в Yandex Cloud

Этот гайд поможет быстро настроить автоматический деплой бота в Yandex Cloud Serverless Containers.

## 🚀 Автоматическая настройка (Рекомендуется)

### 1. Установка Yandex Cloud CLI

```bash
# macOS/Linux
curl -sSL https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash

# Перезапустите терминал или выполните:
source ~/.bashrc
```

### 2. Инициализация Yandex Cloud

```bash
yc init
```

Следуйте инструкциям для:
- Входа в аккаунт
- Выбора облака
- Выбора папки (folder)

### 3. Автоматическая настройка ресурсов

```bash
# Запустите скрипт автоматической настройки
./scripts/setup-yandex-cloud.sh
```

Скрипт автоматически создаст:
- ✅ Сервисный аккаунт с необходимыми ролями
- ✅ Ключ для сервисного аккаунта
- ✅ Container Registry
- ✅ Serverless Container
- ✅ Файл `.env.example` с примером конфигурации

### 4. Настройка GitHub Secrets

Скопируйте вывод скрипта и добавьте секреты в GitHub:

1. Перейдите в **Settings** → **Secrets and variables** → **Actions**
2. Добавьте следующие секреты:

```
YC_CLOUD_ID=<значение_из_скрипта>
YC_FOLDER_ID=<значение_из_скрипта>
YC_REGISTRY_ID=<значение_из_скрипта>
YC_SERVICE_ACCOUNT_ID=<значение_из_скрипта>
YC_SERVICE_ACCOUNT_KEY=<значение_из_скрипта>
YC_CONTAINER_NAME=workouts-bot
BOT_TOKEN=<ваш_токен_telegram_бота>
DATABASE_URL=<url_вашей_базы_данных>
```

### 5. Запуск деплоя

```bash
# Сделайте коммит и пуш в main ветку
git add .
git commit -m "Add Yandex Cloud deployment configuration"
git push origin main
```

Деплой запустится автоматически! 🎉

---

## 🛠 Ручная настройка

Если предпочитаете настроить все вручную, следуйте [подробной инструкции](DEPLOYMENT.md).

---

## 📋 Проверка деплоя

### 1. Мониторинг в GitHub Actions

1. Перейдите в **Actions** в вашем репозитории
2. Найдите workflow **"Deploy to Yandex Cloud"**
3. Проверьте статус выполнения

### 2. Проверка в Yandex Cloud Console

1. Откройте [Yandex Cloud Console](https://console.cloud.yandex.ru/)
2. Перейдите в **Serverless Containers**
3. Найдите контейнер **workouts-bot**
4. Проверьте статус и логи

### 3. Тестирование бота

```bash
# Получите URL контейнера
yc serverless container get workouts-bot --format json | jq -r '.status.url'

# Проверьте health check
curl https://your-container-url.yandexcloud.net/health
```

---

## 🔧 Локальная разработка

### 1. Настройка окружения

```bash
# Скопируйте пример конфигурации
cp .env.example .env

# Отредактируйте .env файл
nano .env
```

### 2. Запуск локально

```bash
# Установка зависимостей
make install

# Запуск в режиме разработки
make dev

# Или обычный запуск
make run
```

### 3. Сборка и тестирование Docker образа

```bash
# Сборка образа
make docker-build

# Тестирование образа
docker run --rm \
  -e BOT_TOKEN=your_token \
  -e DATABASE_URL=your_db_url \
  workouts-bot
```

---

## 🚨 Troubleshooting

### Частые проблемы

**1. Ошибка аутентификации в GitHub Actions**
```
Error: authentication failed
```
**Решение:** Проверьте корректность `YC_SERVICE_ACCOUNT_KEY` в GitHub Secrets

**2. Контейнер не запускается**
```
Error: container failed to start
```
**Решение:**
- Проверьте `BOT_TOKEN` и `DATABASE_URL`
- Посмотрите логи: `yc serverless container revision logs <revision-id>`

**3. Health check не проходит**
```
Health check failed
```
**Решение:**
- Убедитесь, что контейнер отвечает на порту 8080
- Проверьте переменную окружения `PORT`

### Полезные команды

```bash
# Просмотр логов контейнера
yc serverless container revision logs $(yc serverless container revision list --container-name workouts-bot --limit 1 --format json | jq -r '.[0].id')

# Список ревизий
yc serverless container revision list --container-name workouts-bot

# Информация о контейнере
yc serverless container get workouts-bot

# Ручной деплой новой ревизии
yc serverless container revision deploy \
  --container-name workouts-bot \
  --image cr.yandex/$REGISTRY_ID/workouts-bot:latest \
  --cores 1 \
  --memory 512MB \
  --environment BOT_TOKEN=$BOT_TOKEN \
  --environment DATABASE_URL=$DATABASE_URL
```

---

## 📚 Дополнительные ресурсы

- [Подробная документация по деплою](DEPLOYMENT.md)
- [Yandex Cloud Serverless Containers](https://cloud.yandex.ru/docs/serverless-containers/)
- [Container Registry](https://cloud.yandex.ru/docs/container-registry/)
- [GitHub Actions](https://docs.github.com/en/actions)

---

## 🎯 Что дальше?

После успешного деплоя рекомендуется:

1. **Настроить мониторинг** - добавить алерты на ошибки
2. **Настроить базу данных** - использовать Yandex Managed PostgreSQL
3. **Добавить домен** - настроить custom domain для webhook
4. **Настроить CI/CD** - добавить автотесты и staging окружение

Удачного деплоя! 🚀
