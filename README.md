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

## gRPC

| Метод | Описание | Запрос | Ответ |
|------|---------|--------|-------|
| `GetAll` | Получить все события | `EmptyRequest` | `GetAllResponse` |
| `GetAllByCreator` | Получить все события по создателю (UUID) | `GetAllByCreatorRequest` | `GetAllResponse` |
| `GetAllByStatus` | Получить все события по статусу | `GetAllByStatusRequest` | `GetAllResponse` |
| `GetById` | Получить событие по ID | `GetByIdRequest` | `GetByIdResponse` |
| `Create` | Создать новое событие | `CreateRequest` | `CreateResponse` |
| `DeleteById` | Удалить событие по ID | `DeleteByIdRequest` | `DeleteByIdResponse` |
| `Update` | Обновить событие | `UpdateRequest` | `EmptyResponse` |
| `Register` | Зарегистрироваться на событие | `RegisterRequest` | `EmptyResponse` |
| `CancellRegister` | Отменить регистрацию на событие | `CancellRegisterRequest` | `EmptyResponse` |
| `GetAllByUser` | Получить все события, на которые зарегистрирован пользователь | `GetAllByUserRequest` | `GetAllByUserResponse` |
| `GetAllUsersByEvent` | Получить всех пользователей, зарегистрированных на событие | `GetAllUsersByEventRequest` | `GetAllUsersByEventResponse` |

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