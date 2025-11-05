# EventService — gRPC-сервис для управления событиями

Простой и масштабируемый **gRPC-сервис** для управления событиями, написанный на **Go** с использованием **чистой архитектуры**. Поддерживает CRUD-операции, фильтрацию по создателю и статусу, валидацию данных, кэширование через **Redis** и покрытие тестами.

---

## Особенности

- **gRPC API** с Protobuf-определениями
- **Чистая архитектура** (Clean Architecture)
- **Кэширование** через Redis метода GetById
- **Валидация** входных данных 
- **Обработка ошибок** с gRPC-статусами
- **Юнит-тесты** и **интеграционные тесты**
- **Поддержка контекста** и отмены запросов

---

## API (gRPC)

| Метод | Описание | Запрос | Ответ |
|------|---------|-------|-------|
| `GetAll` | Все события | `EmptyRequest` | `GetAllResponse` |
| `GetById` | По ID | `GetByIdRequest` | `GetByIdResponse` |
| `GetAllByCreator` | По создателю (UUID) | `GetAllByCreatorRequest` | `GetAllResponse` |
| `GetAllByStatus` | По статусу | `GetAllByStatusRequest` | `GetAllResponse` |
| `Create` | Создать | `CreateRequest` | `CreateResponse` |
| `Update` | Обновить | `UpdateRequest` | `EmptyResponse` |
| `DeleteById` | Удалить | `DeleteByIdRequest` | `DeleteByIdResponse` |
---

## Шаги по запуску
1. **Клонируй репозиторий и перейдите в папку**:
   ```
   git clone https://github.com/Estriper0/AuthService.git
   cd AuthService
   ```
2. **Настройте переменные окружения в `.env`**:
   ```env
    APP_ENV=local

    DB_HOST=postgres
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=12345
    DB_NAME=event_db
    
    REDIS_PASSWORD=12345
    REDIS_ADDR=redis:6379
   ```
3. **Запусти с помощью Docker Compose**:
   ```
   docker compose up --build -d
   ```

---

## Тестирование

#### Все тесты
```bash
go test ./... -v
```

#### Только юнит-тесты
```bash
go test -short ./... -v
```

> Используется `testcontainers-go` для запуска БД в Docker

---

## Валидация

- Все поля в `CreateRequest` и `UpdateRequest` валидируются

---