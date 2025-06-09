# event-analytics-service

event-analytics-service — это простой сервис для приёма, асинхронной обработки и хранения событий
с мониторингом через Prometheus.

## В проекте реализовано
- Приём ивентов через POST запрос на ручку /events
- Публикация событий в RabbitMQ
- Асинхронная обработка воркерами с сохранением в PostgreSQL
- Метрики статуса обработки и времени работы с Prometheus
- Логгирование ошибок

---

## Технологии

- Go 1.24.2
- RabbitMQ
- PostgreSQL
- Prometheus (метрики)
- Docker и Docker Compose
---