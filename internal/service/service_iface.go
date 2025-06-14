package service

import (
	"context"
	"github.com/diemensa/event-analytics-service/internal/model"
)

type Service interface {
	SaveEvent(ctx context.Context, e *model.Event) error
	GetEvents(ctx context.Context) ([]model.Event, error)
	SendToRabbit(ctx context.Context, e *model.Event) error
}
