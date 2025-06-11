package model

import "time"

type Event struct {
	ID        string    `json:"id" validate:"required,uuid4"`
	Type      string    `json:"type" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	UserID    *string   `json:"user_id,omitempty" validate:"omitempty,uuid4"`
	// указатель для того, чтобы при отсутствии поля не возникало ошибки; будет создаваться NULL в бд
}
