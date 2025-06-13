package repository

import (
	"context"
	"fmt"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestEventRepo_Save_Success(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := NewEventRepo(mockDB)
	userid := uuid.New().String()

	testEvent := &model.Event{
		ID:        uuid.New().String(),
		Type:      "authentication_idk",
		Timestamp: time.Now().UTC(),
		UserID:    &userid,
	}

	query := `INSERT INTO events (id, type, timestamp, user_id) VALUES ($1, $2, $3, $4)`

	mockDB.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(testEvent.ID, testEvent.Type, testEvent.Timestamp, testEvent.UserID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Save(ctx, testEvent)

	assert.NoError(t, err)
	fmt.Println(mockDB.ExpectationsWereMet())
	assert.NoError(t, mockDB.ExpectationsWereMet())

}

func TestEventRepo_Save_DuplicatedErr(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := NewEventRepo(mockDB)
	userid := uuid.New().String()

	testEvent := &model.Event{
		ID:        uuid.New().String(),
		Type:      "authentication_idk",
		Timestamp: time.Now().UTC(),
		UserID:    &userid,
	}

	query := `INSERT INTO events (id, type, timestamp, user_id) VALUES ($1, $2, $3, $4)`

	pgErr := &pgconn.PgError{Code: "23505"}
	mockDB.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(testEvent.ID, testEvent.Type, testEvent.Timestamp, testEvent.UserID).
		WillReturnError(pgErr)

	err = repo.Save(ctx, testEvent)

	assert.ErrorIs(t, err, ErrDuplicatedKey)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestEventRepo_GetEvents(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := NewEventRepo(mockDB)
	userid := uuid.New().String()

	testEvent := model.Event{
		ID:        uuid.New().String(),
		Type:      "authentication_idk",
		Timestamp: time.Now().UTC(),
		UserID:    &userid,
	}

	want := []model.Event{testEvent}

	query := `SELECT id, type, timestamp, user_id FROM events`

	rows := pgxmock.NewRows([]string{"id", "type", "timestamp", "user_id"}).
		AddRow(testEvent.ID, testEvent.Type, testEvent.Timestamp, testEvent.UserID)

	mockDB.ExpectQuery(query).WillReturnRows(rows)

	got, err := repo.GetEvents(ctx)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
