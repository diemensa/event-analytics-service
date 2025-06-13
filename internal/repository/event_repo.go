package repository

import (
	"context"
	"errors"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type EventRepo struct {
	db PGXPool
}

func NewEventRepo(db PGXPool) *EventRepo {
	return &EventRepo{db: db}
}

var ErrDuplicatedKey = errors.New("duplicate key")

func isDuplicateKey(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *EventRepo) Save(ctx context.Context, event *model.Event) error {
	query := `INSERT INTO events (id, type, timestamp, user_id)
VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(ctx, query, event.ID, event.Type, event.Timestamp, event.UserID)

	if err != nil {
		if isDuplicateKey(err) {
			return ErrDuplicatedKey
		}
		return err
	}

	return nil
}

func (r *EventRepo) GetEvents(ctx context.Context) ([]model.Event, error) {
	query := `SELECT id, type, timestamp, user_id FROM events`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var e model.Event
		if err = rows.Scan(&e.ID, &e.Type, &e.Timestamp, &e.UserID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil

}
