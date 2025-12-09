# URL Shortener Service

Сервис для сокращения URL-адресов.

## Функциональность

### Основные возможности:
- **Сокращение URL** через различные интерфейсы:
  - `POST /` - текстовый формат
  - `POST /api/shorten` - JSON формат
  - `POST /api/shorten/batch` - пакетное сокращение URL
- **Перенаправление** по коротким ссылкам: `GET /{id}`
- **Управление ссылками пользователя**:
  - `GET /api/user/urls` - получение всех сокращённых URL пользователя
  - `DELETE /api/user/urls` - асинхронное удаление URL
- **Статистика**: `GET /api/internal/stats` (только для доверенных подсетей)
- **Проверка соединения с БД**: `GET /ping`
- **Поддержка gRPC** - все операции доступны также через gRPC

### Технические особенности:
- Конфигурация через флаги, переменные окружения и JSON-файлы
- Поддержка нескольких хранилищ:
  - PostgreSQL (основное)
  - Файловое хранилище (JSON)
  - In-memory хранилище
- Сжатие данных (gzip) для запросов и ответов
- Аутентификация пользователей через подписанные куки
- Логирование запросов и ответов
- Graceful shutdown
- Поддержка HTTPS
- Статический анализ кода

## Быстрый старт

### Сборка и запуск:

```bash
# Сборка
go build -o shortener ./cmd/shortener

# Запуск с конфигурацией по умолчанию
./shortener

# Запуск с кастомным адресом
./shortener -a localhost:8080 -b http://localhost:8080
```

### Примеры использования:

```bash
# Сокращение URL
curl -X POST -d "https://example.com" http://localhost:8080/

# Сокращение через JSON API
curl -X POST -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}' \
  http://localhost:8080/api/shorten

# Получение оригинального URL
curl -v http://localhost:8080/{short_id}

# Пакетное сокращение
curl -X POST -H "Content-Type: application/json" \
  -d '[{"correlation_id":"1","original_url":"https://example1.com"}]' \
  http://localhost:8080/api/shorten/batch
```

## Конфигурация

Сервис поддерживает несколько способов конфигурации (в порядке приоритета):

1. **Флаги командной строки** (высший приоритет)
2. **Переменные окружения**
3. **JSON-файл конфигурации**
4. **Значения по умолчанию** (низший приоритет)

### Основные параметры:

**Флаги командной строки:**

- `-a` - адрес запуска HTTP-сервера (по умолчанию: `localhost:8080`)
- `-b` - базовый адрес результирующего сокращённого URL
- `-f` - полное имя файла для хранения данных в формате JSON (по умолчанию: `/tmp/short-url-db.json`)
- `-d` - строка подключения к базе данных PostgreSQL
- `-s` - включение HTTPS (true/false)
- `-t` - доверенная подсеть в формате CIDR для доступа к статистике
- `-c` / `-config` - путь к JSON-файлу конфигурации

**Переменные окружения:**

- `SERVER_ADDRESS` - аналог флага `-a`
- `BASE_URL` - аналог флага `-b`
- `FILE_STORAGE_PATH` - аналог флага `-f`
- `DATABASE_DSN` - аналог флага `-d`
- `ENABLE_HTTPS` - аналог флага `-s`
- `TRUSTED_SUBNET` - аналог флага `-t`
- `CONFIG` - аналог флага `-c`

**Пример JSON-конфигурации:**

```json
{
  "server_address": "localhost:8080",
  "base_url": "http://localhost:8080",
  "file_storage_path": "/tmp/short-url-db.json",
  "database_dsn": "postgres://user:pass@localhost:5432/db",
  "enable_https": false,
  "trusted_subnet": "192.168.1.0/24"
}
```

## Лицензия

MIT License
