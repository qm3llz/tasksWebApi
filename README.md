# tasksWebApi

REST API для управления личными задачами: регистрация и авторизация через JWT, CRUD задач с разделением по пользователям, миграции БД.

## Стек

- **Go** + [chi](https://github.com/go-chi/chi) — роутер
- **PostgreSQL** + [pgx](https://github.com/jackc/pgx) — драйвер БД
- **JWT** ([golang-jwt](https://github.com/golang-jwt/jwt)) — авторизация
- **bcrypt** — хеширование паролей
- [golang-migrate](https://github.com/golang-migrate/migrate) — миграции БД
- **Docker** — контейнеризация

## Возможности

- Регистрация и вход (`/register`, `/login`), пароли хешируются через bcrypt
- Все `/tasks`-роуты защищены JWT-аутентификацией
- Задачи изолированы по пользователю на уровне SQL — доступ к чужим задачам невозможен
- Юнит-тесты хендлеров на моках (без реальной БД)

## Быстрый старт

### 1. Настройка окружения

```bash
cp .env.example .env
```

Заполни `.env`:
- `DATABASE_URL` — строка подключения к PostgreSQL
- `JWT_SECRET` — секрет для подписи токенов, сгенерировать:
  ```bash
  openssl rand -hex 32
  ```

### 2. Миграции

```bash
migrate -path ./migrations -database "$DATABASE_URL" up
```

### 3. Запуск

```bash
go run ./cmd/api
```

Сервер стартует на `:8080`.

### Через Docker

```bash
docker build -t tasks-api .
docker run --env-file .env -p 8080:8080 tasks-api
```

## API

### Открытые роуты

| Метод | Путь | Описание |
|---|---|---|
| POST | `/register` | Регистрация (`username`, `password`) |
| POST | `/login` | Вход, возвращает JWT |

### Защищённые роуты (требуют `Authorization: Bearer <token>`)

| Метод | Путь | Описание |
|---|---|---|
| POST | `/tasks` | Создать задачу |
| GET | `/tasks` | Список задач текущего пользователя |
| GET | `/tasks/{id}` | Получить задачу по id |
| PUT | `/tasks/{id}` | Обновить задачу |
| DELETE | `/tasks/{id}` | Удалить задачу |

## Модель задачи

```go
type Task struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    Name        string
    Status      TaskStatus // "todo" | "in_progress" | "done"
    Description *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Тесты

```bash
go test ./... -v
```

Хендлеры покрыты табличными юнит-тестами через фейковый репозиторий (`fakeRepo`), реальная БД для тестов не требуется.

## Структура проекта

```
cmd/api/            точка входа
internal/
  auth/              генерация и проверка JWT
  db/                подключение к PostgreSQL
  handler/           HTTP-хендлеры (tasks, users)
  middleware/        AuthMiddleware
  models/            доменные модели
  repository/        доступ к БД (tasks, users)
migrations/          SQL-миграции (golang-migrate)
```
