package worker

import (
	"context"
	"encoding/json"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/diemensa/event-analytics-service/internal/repository"
	"github.com/diemensa/event-analytics-service/internal/service"
	"log"

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
			for msg := range p.events {
				var event model.Event
				err := json.Unmarshal(msg.Body, &event)

				log.Printf("worker %d is processing event %s", workerID, event.ID)

				if err != nil {
					log.Printf("worker %d failed to unmarshal event: %v", workerID, err)
					_ = msg.Nack(false, false)
					continue
				}

				err = p.service.SaveEvent(ctx, &event)
				if err != nil {
					log.Printf("worker %d failed to save event in DB: %v", workerID, err)
					_ = msg.Nack(false, true)
					continue
				}

				_ = msg.Ack(false)
				log.Printf("worker %d successfully finished event %v", workerID, event.ID)

			}
		}(i)
	}
}
