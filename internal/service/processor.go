package service

import (
	"context"
	"github.com/diemensa/event-analytics-service/internal/broker"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/diemensa/event-analytics-service/internal/repository"
)

type EventService struct {
	rabbitRepo broker.Publisher
	eventRepo  *repository.EventRepo
}

func NewEventService(
	r broker.Publisher,
	e *repository.EventRepo) *EventService {
	return &EventService{
		rabbitRepo: r,
		eventRepo:  e,
	}
}

func (service *EventService) SaveEvent(ctx context.Context, e *model.Event) error {
	return service.eventRepo.Save(ctx, e)
}

func (service *EventService) GetEvents(ctx context.Context) ([]model.Event, error) {
	return service.eventRepo.GetEvents(ctx)
}

func (service *EventService) SendToRabbit(ctx context.Context, e *model.Event) error {
	return service.rabbitRepo.Publish(ctx, e)
}
