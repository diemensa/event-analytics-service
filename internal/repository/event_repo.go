package repository

import (
	"context"
	"github.com/diemensa/event-analytics-service/internal/model"
	"gorm.io/gorm"
	"strings"
)

type EventRepo struct {
	db *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Save(ctx context.Context, event *model.Event) error {
	err := r.db.WithContext(ctx).Create(event).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(),
			"unique constraint") {
			return gorm.ErrDuplicatedKey
		}
	}
	return nil
}

func (r *EventRepo) GetEvents(ctx context.Context) ([]model.Event, error) {
	var events []model.Event

	err := r.db.WithContext(ctx).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil

}
