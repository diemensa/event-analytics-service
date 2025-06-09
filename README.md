# event-analytics-service

event-analytics-service — это микросервис для приёма, асинхронной обработки и хранения событий
с мониторингом через Prometheus и визуализацией через Grafana

## В проекте реализовано

- Приём ивентов через POST запрос на ручку /events
- Публикация событий в очереди RabbitMQ
- Асинхронная обработка воркерами с сохранением в PostgreSQL
- Метрики статуса обработки и времени работы
- Логирование ошибок

---

## Технологии
- Go 1.24.2
- GORM
- RabbitMQ
- PostgreSQL
- Prometheus (метрики посылаются на localhost:8080/metrics)
- Grafana (дашбордики создаются на localhost:53000, admin/admin)
- Docker
---