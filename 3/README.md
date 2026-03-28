# Calendar HTTP Service

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
![HTTP](https://img.shields.io/badge/HTTP-Service-0A0A0A)
![Tests](https://img.shields.io/badge/tests-unit%20%7C%20handler-green)
![Logger](https://img.shields.io/badge/logger-async-blue)
![Workers](https://img.shields.io/badge/workers-reminder%20%7C%20archive-orange)
![Storage](https://img.shields.io/badge/storage-in--memory-lightgrey)

Полноценный HTTP-сервис календаря на Go с CRUD-операциями, асинхронным логированием, фоновым воркером напоминаний и воркером архивации старых событий.

---

## Возможности

- создание события
- обновление события
- удаление события
- получение событий на день
- получение событий на неделю
- получение событий на месяц
- асинхронный логгер через канал
- фоновый reminder worker через канал
- фоновый archive worker по таймеру
- unit- и handler-тесты

---

## Архитектура

Проект разделен на несколько слоев:

- `cmd/main.go` — точка входа и сборка зависимостей
- `internal/handlers` — HTTP-хендлеры
- `internal/middleware` — middleware логирования
- `internal/service` — бизнес-логика
- `internal/repository` — in-memory хранилище
- `internal/worker` — фоновые воркеры
- `internal/myLogger` — асинхронный логгер
- `internal/dto` — входные DTO
- `internal/customErrs` — доменные ошибки

---

## Как запустить

### 1. Установить зависимости

```bash
go mod tidy
```

### 2. Запустить сервер

```bash
go run ./cmd -port=8080
```

Сервер будет доступен по адресу:

```text
http://localhost:8080
```

---

## Как запустить тесты

```bash
go test ./...
```

---

## Используемый формат дат

Во всех запросах даты передаются в формате **RFC3339**.

Пример:

```text
2026-03-27T23:40:00+03:00
```

Важно:

- для `POST`-ручек данные передаются в `x-www-form-urlencoded`
- для `GET`-ручек данные передаются через query params
- в query string символ `+` лучше передавать как `%2B` или использовать вкладку **Params** в Postman

---

## Endpoints

### `GET /status`

Проверка, что сервер запущен.

#### Response

```json
{
  "result": "server is running"
}
```

---

### `POST /create_event`

Создание нового события.

#### Body (`x-www-form-urlencoded`)

- `user_id` — идентификатор пользователя
- `title` — название события
- `date` — время события в RFC3339
- `description` — описание события
- `priority` — приоритет (`high`, `medium`, `low`)
- `remind_at` — необязательное время напоминания в RFC3339

#### Пример

- `user_id=1`
- `title=team meeting`
- `date=2026-03-27T23:50:00+03:00`
- `description=weekly sync`
- `priority=high`
- `remind_at=2026-03-27T23:45:00+03:00`

#### Response

```json
{
  "id": 1,
  "result": "event created"
}
```

---

### `GET /events_for_day`

Получение событий пользователя за день.

#### Query Params

- `user_id`
- `date`

#### Пример

```text
/events_for_day?user_id=1&date=2026-03-27T00:00:00%2B03:00
```

#### Response

```json
{
  "result": [
    {
      "id": 1,
      "user_id": 1,
      "title": "team meeting"
    }
  ]
}
```

---

### `GET /events_for_week`

Получение событий пользователя за неделю.

#### Query Params

- `user_id`
- `date`

#### Пример

```text
/events_for_week?user_id=1&date=2026-03-27T00:00:00%2B03:00
```

---

### `GET /events_for_month`

Получение событий пользователя за месяц.

#### Query Params

- `user_id`
- `date`

#### Пример

```text
/events_for_month?user_id=1&date=2026-03-27T00:00:00%2B03:00
```

---

### `POST /update_event`

Обновление существующего события.

#### Body (`x-www-form-urlencoded`)

- `id`
- `title`
- `date`
- `description`
- `priority`
- `remind_at`

#### Пример

- `id=1`
- `title=updated meeting`
- `date=2026-03-28T00:10:00+03:00`
- `description=updated description`
- `priority=medium`
- `remind_at=2026-03-28T00:05:00+03:00`

#### Response

```json
{
  "result": "event updated"
}
```

---

### `POST /delete_event`

Удаление события.

#### Body (`x-www-form-urlencoded`)

- `id`

#### Пример

- `id=1`

#### Response

```json
{
  "result": "event deleted"
}
```

---

## Примеры запросов

### cURL: create event

```bash
curl -X POST "http://localhost:8080/create_event" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_id=1" \
  -d "title=team meeting" \
  -d "date=2026-03-27T23:50:00+03:00" \
  -d "description=weekly sync" \
  -d "priority=high" \
  -d "remind_at=2026-03-27T23:45:00+03:00"
```

### cURL: get events for day

```bash
curl "http://localhost:8080/events_for_day?user_id=1&date=2026-03-27T00:00:00%2B03:00"
```

### cURL: update event

```bash
curl -X POST "http://localhost:8080/update_event" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=1" \
  -d "title=updated meeting" \
  -d "date=2026-03-28T00:10:00+03:00" \
  -d "description=updated description" \
  -d "priority=medium" \
  -d "remind_at=2026-03-28T00:05:00+03:00"
```

### cURL: delete event

```bash
curl -X POST "http://localhost:8080/delete_event" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=1"
```

---

## Фоновые воркеры

### Reminder Worker

Если событие создается с `remind_at`, после успешного создания в канал помещается задача `ReminderTask`.

Отдельная горутина:

- читает задачи из канала
- вычисляет время ожидания до `remind_at`
- по наступлении времени пишет лог о сработавшем напоминании

Пример лога:

```text
2026-03-27T23:40:00+03:00 [INFO] reminder sent for event id=2 title=test1
```

Если время напоминания уже прошло, задача пропускается:

```text
2026-03-27T23:39:10+03:00 [INFO] reminder skipped 2 test1 because remind time already passed
```

---

### Archive Worker

Отдельная горутина по таймеру:

- получает все активные события
- ищет события, время которых уже прошло
- переносит их в архив

Пример лога:

```text
2026-03-27T23:45:00+03:00 [INFO] event moved to archive
```

Если архивирование завершилось с ошибкой:

```text
2026-03-27T23:45:00+03:00 [ERROR] failed move to archive err=...
```

---

## Асинхронный логгер

Все логи пишутся не напрямую из HTTP-хендлеров, а через буферизированный канал.

Логгер:

- принимает структурированные сообщения
- обрабатывает их в отдельной горутине
- выводит в едином формате

Логируются:

- HTTP-запросы
- ошибки сервиса
- ошибки сериализации ответа
- срабатывание reminder worker
- архивация событий

Пример access log:

```text
2026-03-27T23:38:31+03:00 [INFO] request handled method=POST path=/create_event status=200 duration=1.2ms
```

---

## Ручной сценарий проверки

### Проверка `status`

Открыть:

```text
GET http://localhost:8080/status
```

Ожидаемый ответ:

```json
{
  "result": "server is running"
}
```

---

### Проверка reminder worker

1. Запустить сервер
2. Создать событие с `remind_at` через 30–60 секунд
3. Убедиться, что сервер вернул `event created`
4. Дождаться времени напоминания
5. Убедиться, что в логах появилось сообщение от reminder worker

---

### Проверка archive worker

1. Создать событие с `date`, которая наступит в ближайшее время
2. Дождаться, пока `eventTime` станет прошлым
3. Дождаться следующего запуска archive worker
4. Убедиться, что в логах появилось сообщение о переносе события в архив

---

## Покрытие тестами

В проекте реализованы тесты для:

- `service.Create`
- создания события с напоминанием и без него
- ошибок валидации времени
- HTTP-хендлеров:
    - `status`
    - `create_event`
    - `events_for_day`
    - `events_for_week`
    - `events_for_month`
    - `delete_event`
    - `update_event`
- repository-логики переноса в архив

Запуск всех тестов:

```bash
go test ./...
```

---

## Примечания

- хранилище in-memory, данные не сохраняются между перезапусками сервера
- даты в query-параметрах должны передаваться как корректный RFC3339
- для временной ручной проверки archive worker можно уменьшить интервал, затем вернуть значение обратно
