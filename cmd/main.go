package main

import (
	"context"
	"github.com/diemensa/event-analytics-service/config"
	"github.com/diemensa/event-analytics-service/internal/broker"
	"github.com/diemensa/event-analytics-service/internal/handler"
	"github.com/diemensa/event-analytics-service/internal/metrics"
	"github.com/diemensa/event-analytics-service/internal/repository"
	"github.com/diemensa/event-analytics-service/internal/service"
	"github.com/diemensa/event-analytics-service/internal/worker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadEnv()

	db, err := config.InitPostgres(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal("couldn't connect to database:", err.Error())
	}

	const workerCount = 10
	ctx := context.Background()

	rabbit, err := broker.NewRabbitPublisher(cfg.RabbitMQURI, cfg.RabbitQueueName)
	if err != nil {
		log.Fatalf("couldn't connect to rabbitMQ: %v", err)
	}
	defer func() {
		if err = rabbit.Close(); err != nil {
			log.Printf("failed to close rabbit publisher: %v", err)
		}
	}()

	eventRepo := repository.NewEventRepo(db)
	eventService := service.NewEventService(rabbit, eventRepo)

	messageChan, err := rabbit.Consume(cfg.RabbitQueueName)
	if err != nil {
		log.Fatalf("failed to start consuming rabbit channel: %v", err)
	}
	workerPool := worker.NewPool(eventService, eventRepo, messageChan)
	workerPool.Start(ctx, workerCount)

	eventHandler := handler.NewEventHandler(eventService)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			eventHandler.NewGetEventHandler(w, r)
		case http.MethodPost:
			eventHandler.HandleCreateEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	metrics.PrometheusInit()
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))

}
