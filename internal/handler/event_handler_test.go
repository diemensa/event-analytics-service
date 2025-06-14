package handler

import (
	"bytes"
	"encoding/json"
	"github.com/diemensa/event-analytics-service/internal/handler/dto"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/diemensa/event-analytics-service/internal/repository"
	"github.com/diemensa/event-analytics-service/internal/service"
	"github.com/diemensa/event-analytics-service/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEventHandler_HandleCreateEvent_Accepted(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	evRepo := repository.NewEventRepo(mockDB)
	rabbitRepo := mocks.NewPublisher(t)
	serv := service.NewEventService(rabbitRepo, evRepo)

	userid := uuid.New().String()
	testEvent := model.Event{
		ID:        uuid.New().String(),
		Type:      "authentication_idk",
		Timestamp: time.Now().UTC(),
		UserID:    &userid,
	}

	handler := NewEventHandler(serv)

	body, err := json.Marshal(testEvent)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	defer req.Body.Close()

	publishCalled := make(chan struct{})

	rabbitRepo.On("Publish", mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			close(publishCalled)
		}).Once() // мощный костыль в тестах > создание waitgroup в структуре

	handler.HandleCreateEvent(rec, req)

	<-publishCalled

	res := rec.Result()

	respBody := dto.CreateEventResponse{
		ID:     testEvent.ID,
		Status: "accepted",
	}

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	assert.Equal(t, testEvent.ID, respBody.ID)
	assert.Equal(t, "accepted", respBody.Status)
}

func TestEventHandler_HandleCreateEvent_InvalidRequest(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	evRepo := repository.NewEventRepo(mockDB)
	rabbitRepo := mocks.NewPublisher(t)
	serv := service.NewEventService(rabbitRepo, evRepo)

	handler := NewEventHandler(serv)

	body := `{"id": "123", "typ": "idk", "timestamp": "whatever"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	defer req.Body.Close()

	handler.HandleCreateEvent(rec, req)

	res := rec.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	respBody, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(respBody), "invalid request body")
}

func TestEventHandler_HandleCreateEvent_ValidationFail(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	evRepo := repository.NewEventRepo(mockDB)
	rabbitRepo := mocks.NewPublisher(t)
	serv := service.NewEventService(rabbitRepo, evRepo)

	handler := NewEventHandler(serv)

	body := `{"id": "00000000-0000-0000-0000-000000000000", "type": "test", "timestamp": "2025-06-14T10:00:00Z"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	defer req.Body.Close()

	handler.HandleCreateEvent(rec, req)

	res := rec.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	respBody, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(respBody), "validation error")
}

func TestEventHandler_HandleGetEvents(t *testing.T) {

	serv := mocks.NewService(t)

	handler := NewEventHandler(serv)

	userid := uuid.New().String()
	testEvent := model.Event{
		ID:        uuid.New().String(),
		Type:      "authentication_idk",
		Timestamp: time.Now().UTC(),
		UserID:    &userid,
	}

	want := []model.Event{testEvent}
	serv.On("GetEvents", mock.Anything).Return(want, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	defer req.Body.Close()

	handler.HandleGetEvents(rec, req)

	res := rec.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	respBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var got []model.Event
	err = json.Unmarshal(respBody, &got)
	assert.NoError(t, err)

	assert.Equal(t, want, got)
	serv.AssertExpectations(t)
}
