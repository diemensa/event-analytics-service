package repository

import (
	"context"
	"github.com/diemensa/event-analytics-service/internal/model"
	"gorm.io/gorm"
)

type EventRepo struct {
	db *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Save(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *EventRepo) GetEvents(ctx context.Context) ([]model.Event, error) {
	var events []model.Event

	err := r.db.WithContext(ctx).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil

}
