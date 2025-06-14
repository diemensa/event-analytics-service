package handler

import (
	"context"
	"encoding/json"
	"github.com/diemensa/event-analytics-service/internal/handler/dto"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/diemensa/event-analytics-service/internal/service"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type EventHandler struct {
	eventService service.Service
}

func NewEventHandler(s service.Service) *EventHandler {
	return &EventHandler{eventService: s}
}

func (h *EventHandler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {

	var validate = validator.New()
	var event model.Event
	decoder := json.NewDecoder(r.Body)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}()

	if err := decoder.Decode(&event); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(event); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	go func(ctx context.Context, e model.Event) {
		if err := h.eventService.SendToRabbit(ctx, &e); err != nil {
			log.Printf("failed to send event to RabbitMQ: %v", err)
		}
	}(ctx, event)

	resp := dto.CreateEventResponse{
		ID:     event.ID,
		Status: "accepted",
	}

	err := writeJson(w, http.StatusAccepted, resp)
	if err != nil {
		log.Printf("failed to encode event: %v", err)
	}

}

func (h *EventHandler) HandleGetEvents(w http.ResponseWriter, r *http.Request) {

	events, err := h.eventService.GetEvents(r.Context())

	if err != nil {
		http.Error(w, "failed to query events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = writeJson(w, http.StatusOK, events); err != nil {
		http.Error(w, "failed to encode events: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func writeJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
