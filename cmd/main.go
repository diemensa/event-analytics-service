package main

import (
	"context"
	"errors"
	"github.com/diemensa/event-analytics-service/config"
	"github.com/diemensa/event-analytics-service/internal/broker"
	"github.com/diemensa/event-analytics-service/internal/handler"
	"github.com/diemensa/event-analytics-service/internal/repository"
	"github.com/diemensa/event-analytics-service/internal/service"
	"github.com/diemensa/event-analytics-service/internal/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadEnv()

	db, err := config.InitPostgres(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatal("couldn't connect to database:", err.Error())
	}

	const workerCount = 10
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rabbitRepo, err := broker.NewRabbitPublisher(cfg.RabbitMQURI, cfg.RabbitQueueName)
	if err != nil {
		log.Fatalf("couldn't connect to rabbitMQ: %v", err)
	}
	defer func() {
		if err = rabbitRepo.Close(); err != nil {
			log.Printf("failed to close rabbit publisher: %v", err)
		}
	}()

	eventRepo := repository.NewEventRepo(db)
	eventService := service.NewEventService(rabbitRepo, eventRepo)

	messageChan, err := rabbitRepo.Consume()
	if err != nil {
		log.Fatalf("failed to start consuming rabbit channel: %v", err)
	}
	workerPool := worker.NewPool(eventService, eventRepo, messageChan)
	workerPool.Start(ctx, workerCount)

	handler.SetupHandlers(eventService)

	server := &http.Server{Addr: ":8080", Handler: nil}

	go func() {
		log.Println("server started on :8080")
		if err = server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	log.Println("shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("server couldn't shutdown gracefully: %v", err)
	}

	log.Println("server closed")
}
