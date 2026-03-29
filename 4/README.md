# GC Analyzer

Небольшой HTTP-сервер на Go, который отдает текущую информацию о памяти и сборщике мусора в формате **Prometheus**.

В проекте используются:

- `runtime.ReadMemStats` — для чтения runtime-метрик памяти и GC;
- `debug.SetGCPercent` — для настройки агрессивности сборщика мусора;
- `net/http/pprof` — для профилирования приложения.

## Возможности

Сервер предоставляет:

- endpoint `/metrics` с метриками в формате Prometheus;
- endpoint `/debug/pprof/` и связанные pprof-роуты для профилирования;
- настройку `GCPercent` через флаг запуска.

## Примеры экспортируемых метрик

Сервер отдает, например, такие метрики:

- количество аллокаций;
- количество завершенных циклов GC;
- текущее использование памяти;
- объем heap-памяти;
- количество объектов в heap;
- последнее время запуска GC;
- объем памяти, полученной от ОС.

## Структура проекта

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── httpserver/
│   │   └── router.go
│   └── metrics/
│       └── handler.go
├── go.mod
└── README.md
```

## Запуск

### 1. Перейти в директорию проекта

```bash
cd <директория_проекта>
```

### 2. Запустить сервер

```bash
go run ./cmd/server
```

По умолчанию сервер стартует на:

```text
:8080
```

### 3. Запуск с параметрами

Можно изменить адрес сервера и значение `GCPercent`:

```bash
go run ./cmd/server -addr=:9090 -gc-percent=200
```

Где:

- `-addr` — адрес HTTP-сервера;
- `-gc-percent` — значение для `debug.SetGCPercent`.

## Параметры запуска

| Флаг | По умолчанию | Описание |
|------|--------------|----------|
| `-addr` | `:8080` | адрес HTTP-сервера |
| `-gc-percent` | `100` | целевой процент роста heap до следующего запуска GC |

## HTTP endpoints

### `/`

Проверочный endpoint доступности сервера.

Пример:

```bash
curl http://localhost:8080/
```

Ответ:

```text
server running
```

### `/metrics`

Endpoint с метриками в формате Prometheus.

Пример:

```bash
curl http://localhost:8080/metrics
```

Пример ответа:

```text
# HELP app_up Shows that the application is running
# TYPE app_up gauge
app_up 1

# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 185432

# HELP go_memstats_allocated_bytes_total Total number of bytes allocated even if freed
# TYPE go_memstats_allocated_bytes_total counter
go_memstats_allocated_bytes_total 912384

# HELP go_memstats_mallocs_total Total number of mallocs
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1542

# HELP go_memstats_sys_bytes Number of bytes obtained from system
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 6642704

# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 185432

# HELP go_memstats_heap_objects Number of allocated heap objects
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 927

# HELP go_gc_cycles_total Number of completed GC cycles
# TYPE go_gc_cycles_total counter
go_gc_cycles_total 3

# HELP go_gc_last_time_seconds Unix time of the last GC in seconds
# TYPE go_gc_last_time_seconds gauge
go_gc_last_time_seconds 1711734000.123456789
```

## Профилирование через pprof

В приложении подключены стандартные pprof endpoints:

- `/debug/pprof/`
- `/debug/pprof/heap`
- `/debug/pprof/profile`
- `/debug/pprof/goroutine`
- `/debug/pprof/allocs`
- `/debug/pprof/block`
- `/debug/pprof/mutex`
- `/debug/pprof/threadcreate`
- `/debug/pprof/trace`

### Открыть pprof index в браузере

```text
http://localhost:8080/debug/pprof/
```

### Посмотреть heap profile

```bash
go tool pprof http://localhost:8080/debug/pprof/heap
```

### Посмотреть all allocations profile

```bash
go tool pprof http://localhost:8080/debug/pprof/allocs
```

### Снять CPU profile за 10 секунд

```bash
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=10
```

### Снять trace за 5 секунд

```bash
curl "http://localhost:8080/debug/pprof/trace?seconds=5" -o trace.out
go tool trace trace.out
```

## Какие данные используются

Метрики формируются вручную на основе `runtime.MemStats`, полученного через:

```go
var mem runtime.MemStats
runtime.ReadMemStats(&mem)
```

Используются, например, такие поля:

- `Alloc`
- `TotalAlloc`
- `Mallocs`
- `Sys`
- `HeapAlloc`
- `HeapObjects`
- `NumGC`
- `LastGC`

## Что показывает программа

Программа позволяет наблюдать:

- текущее использование памяти;
- суммарный объем аллокаций;
- количество выделений памяти;
- количество циклов сборки мусора;
- время последнего запуска GC;
- состояние heap.

## Пример проверки работы

После запуска можно выполнить:

```bash
curl http://localhost:8080/
curl http://localhost:8080/metrics
go tool pprof http://localhost:8080/debug/pprof/heap
```
