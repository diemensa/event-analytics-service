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
- Визуализация времени обработки и общего количества ивентов со статусами через Grafana

---

## Технологии

- Go 1.24.2
- RabbitMQ
- PostgreSQL (через pgx)
- Prometheus (метрики посылаются на localhost:8080/metrics)
- Grafana (localhost:3000, admin/admin)
- Docker
---