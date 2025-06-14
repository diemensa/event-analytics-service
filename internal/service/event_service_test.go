package service

import (
	"context"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/diemensa/event-analytics-service/internal/repository"
	"github.com/diemensa/event-analytics-service/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestService_SaveEvent(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	evRepo := repository.NewEventRepo(mockDB)
	rabbitRepo := mocks.NewPublisher(t)
	serv := NewEventService(rabbitRepo, evRepo)

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

	err = serv.SaveEvent(ctx, testEvent)

	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())

}

func TestService_GetEvents(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	evRepo := repository.NewEventRepo(mockDB)
	rabbitRepo := mocks.NewPublisher(t)
	serv := NewEventService(rabbitRepo, evRepo)

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

	got, err := serv.GetEvents(ctx)

	assert.NoError(t, err)
	assert.Equal(t, want, got)

}

func TestService_SendToRabbit(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	evRepo := repository.NewEventRepo(mockDB)
	rabbitRepo := mocks.NewPublisher(t)
	serv := NewEventService(rabbitRepo, evRepo)

	userid := uuid.New().String()
	testEvent := &model.Event{
		ID:        uuid.New().String(),
		Type:      "authentication_idk",
		Timestamp: time.Now().UTC(),
		UserID:    &userid,
	}

	rabbitRepo.On("Publish", ctx, testEvent).Return(nil)

	err = serv.SendToRabbit(ctx, testEvent)

	assert.NoError(t, err)
}
