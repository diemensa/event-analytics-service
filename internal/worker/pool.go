package worker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/diemensa/event-analytics-service/internal/metrics"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/diemensa/event-analytics-service/internal/repository"
	"github.com/diemensa/event-analytics-service/internal/service"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Pool struct {
	service   *service.EventService
	eventRepo *repository.EventRepo
	events    <-chan amqp.Delivery
}

func NewPool(service *service.EventService, repo *repository.EventRepo, events <-chan amqp.Delivery) *Pool {
	return &Pool{
		service:   service,
		eventRepo: repo,
		events:    events,
	}
}

func (p *Pool) Start(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					log.Printf("worker %d stopped by context", workerID)
					return

				case msg, ok := <-p.events:
					if !ok {
						log.Printf("worker %d: events channel closed", workerID)
						return
					}

					startTime := time.Now()

					var event model.Event
					err := json.Unmarshal(msg.Body, &event)

					log.Printf("worker %d is processing event %s", workerID, event.ID)

					if err != nil {
						log.Printf("worker %d failed to unmarshal event: %v", workerID, err)
						metrics.MessagesProcessed.WithLabelValues("fail").Inc()
						_ = msg.Nack(false, false)
						continue
					}

					err = p.service.SaveEvent(ctx, &event)
					if err != nil {
						if errors.Is(err, repository.ErrDuplicatedKey) {
							log.Printf("worker %d: event %s already exists, skipping", workerID, event.ID)
							metrics.MessagesProcessed.WithLabelValues("fail").Inc()
							_ = msg.Nack(false, false)
							continue
						}

						log.Printf("worker %d failed to save event in DB: %v", workerID, err)
						metrics.MessagesProcessed.WithLabelValues("fail").Inc()
						_ = msg.Nack(false, true)
						continue
					}

					_ = msg.Ack(false)
					log.Printf("worker %d successfully finished event %v", workerID, event.ID)
					metrics.MessagesProcessed.WithLabelValues("success").Inc()
					metrics.MessageProcessingDuration.Observe(time.Since(startTime).Seconds())

				}
			}
		}(i)
	}
}
