# event-analytics-service

event-analytics-service — это микросервис для приёма, асинхронной обработки и хранения событий
с мониторингом через Prometheus и визуализацией через Grafana

---

## В проекте реализовано

- Приём ивентов через POST запрос на ручку /events
- Публикация событий в очереди RabbitMQ
- Асинхронная обработка воркерами с сохранением в PostgreSQL
- Логирование ошибок
- Метрики через Prometheus
- Визуализация времени обработки и rate обработанных ивентов со статусами через Grafana
- Юнит-тестирование

---

## Технологии

- Go 1.24.2
- RabbitMQ
- PostgreSQL (через pgx)
- Prometheus (метрики посылаются на localhost:8080/metrics)
- Grafana (localhost:3000, admin/admin)
- Testify + Mockery
- Docker
---

## Запуск

### Linux/MacOS
1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/diemensa/event-analytics-service
   cd event-analytics-service
2. Запустить тесты и приложение:
   ```bash
   make

### Windows
1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/diemensa/event-analytics-service
   cd event-analytics-service
2. Запустить тесты:
   ```bash
   go test -v ./...
3. Собрать и запустить проект через docker-compose:
   ```bash
   docker-compose up --build
