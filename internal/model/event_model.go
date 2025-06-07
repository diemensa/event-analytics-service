package model

import "time"

type Event struct {
	ID        string    `gorm:"primaryKey" json:"id" validate:"required,uuid4"`
	Type      string    `json:"type" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	UserID    string    `json:"user_id,omitempty" validate:"omitempty"`
}
