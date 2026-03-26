# Calendar HTTP Server
Простой HTTP-сервер для управления событиями в календаре.
### Быстрый старт
```
go run main.go -port=8080
```
Если не указывать флаг `port`, то по умолчанию сервер запустится на порту `8080`.

### Проверить работу 
`curl http://localhost:8080/status`
### Доступны следующие эндпоинты
```
POST /create_event — создание нового события
```
```
POST /update_event — обновление существующего
```
```
POST /delete_event — удаление
```
```
GET /events_for_day — получить все события на день
```
```
GET /events_for_week — события на неделю
```
```
GET /events_for_month — события на месяц
```
### Тестирование
```
go test ./... -v
```
```
go test -race ./...
```
```
go vet ./...
```
## Структура проекта
```
2_18/
├── main.go          # HTTP сервер и handlers
├── calendar/
│   ├── calendar.go  # Бизнес-логика календаря
│   └── calendar_test.go
└── models/
└── event.go     # Модель события
```