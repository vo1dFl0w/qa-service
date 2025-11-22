# QA Service

API-сервис «Вопросы — Ответы» (Questions & Answers).
Реализован на Go с использованием `net/http`, GORM и PostgreSQL.

---

## Технологии

* Язык: Go
* База данных: PostgreSQL
* Контейнеризация: Docker + docker-compose
* Миграции: goose

---

---

# Структура проекта (основные папки)

```
cmd/qa-service               # entrypoint
configs/                    # config.yaml
internal/
  app/
    adapters/
      storage/              # Postgres adapters + teststore
    domain/                 # сущности (Question, Answer), ошибки
    repository/             # интерфейсы репозиториев
    usecase/                # бизнес-логика
  transport/
    http_transport/         # handlers, router, middlewares
  config/                   # конфиг-парсинг
  logger/                   # логирование
migrations/                 # SQL миграции
Dockerfile
docker-compose.yaml
```

---

# Быстрый старт (Docker)

Проект уже сконфигурирован для запуска через `docker-compose`.

1. Клонируйте репозиторий:

```bash
git clone https://github.com/vo1dFl0w/qa-service.git
cd qa-service
```

2. Запустите всё в Docker:

```bash
docker-compose up --build
```

`docker-compose` поднимет сервис приложения и PostgreSQL; также при старте выполняются миграции (миграции находятся в `migrations/`).

3. По умолчанию приложение доступно на `http://localhost:8080` (порт configurable через `configs/config.yaml`).

---

# Конфигурация

Вы можете задать/переопределить параметры через правку `configs/config.yaml` (см. `internal/config/config.go`).

---

# Миграции

Миграции SQL находятся в папке `migrations/`:

* `20251120131050_questions.sql`
* `20251120131103_answers.sql`

При запуске через `docker-compose up` миграции выполняются автоматически.

---

## Domain entity `(/internal/domain/)`

Модели:

* **Question**

  * `id: int`
  * `text: string`
  * `created_at: time.Time`
* **Answer**

  * `id: int`
  * `question_id: int` — ссылка на Question
  * `user_id: uuid.UUID` — идентификатор пользователя
  * `text: string`
  * `created_at: time.Time`

API:

* Вопросы:

  * `GET  /questions/` — список всех вопросов
  * `POST /questions/` — создать вопрос
  * `GET  /questions/{id}` — получить вопрос + все ответы на него
  * `DELETE /questions/{id}` — удалить вопрос (каскадно удаляются ответы)
* Ответы:

  * `POST /questions/{id}/answers/` — добавить ответ к вопросу
  * `GET  /answers/{id}` — получить конкретный ответ
  * `DELETE /answers/{id}` — удалить ответ

Бизнес-правила:

* Нельзя создать ответ для несуществующего вопроса (usecase проверяет наличие вопроса).
* Один пользователь может оставлять несколько ответов на один вопрос.
* При удалении вопроса ответы удаляются каскадно (реализовано в миграциях/моделях).

---

# API — примеры запросов (curl)

### Создать вопрос

```bash
curl -X POST http://localhost:8080/questions \
  -H "Content-Type: application/json" \
  -d '{"text": "question_1"}'
```

Ответ: `201 Created`

### Получить список вопросов

```bash
curl http://localhost:8080/questions
```

### Получить вопрос и его ответы

```bash
curl http://localhost:8080/questions/1
```

Ответ содержит вопрос и массив `answers`.

### Удалить вопрос (и все его ответы)

```bash
curl -X DELETE http://localhost:8080/questions/1
```

### Добавить ответ к вопросу

```bash
curl -X POST http://localhost:8080/questions/1/answers \
  -H "Content-Type: application/json" \
  -d '{"user_id":"123e4567-e89b-12d3-a456-426655440000","text":"answer_1"}'
```

`user_id` - можно оставить nil (сделано для упрощения тестирования).

### Получить ответ

```bash
curl http://localhost:8080/answers/1
```

### Удалить ответ

```bash
curl -X DELETE http://localhost:8080/answers/1
```

---

# Тесты

Запускать локально:

```bash
go test ./... -v
```

В проекте присутствуют:

* unit-тесты для usecase (работают с `teststore`),
* unit-тесты для HTTP handlers (используется `httptest` + mocks),
* mocks для репозиториев (в `internal/mocks`).

---
