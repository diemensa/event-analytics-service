package handler

import (
	"context"
	"encoding/json"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/diemensa/event-analytics-service/internal/service"
	"log"
	"net/http"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(s *service.EventService) *EventHandler {
	return &EventHandler{eventService: s}
}

func (h *EventHandler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {

	var event model.Event
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&event); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	go func(e model.Event) {
		ctx := context.Background()
		if err := h.eventService.SendToRabbit(ctx, &e); err != nil {
			log.Printf("failed to send event to RabbitMQ: %v", err)
		}
	}(event)

	w.WriteHeader(http.StatusAccepted)
}

func (h *EventHandler) NewGetEventHandler(w http.ResponseWriter, r *http.Request) {

	events, err := h.eventService.GetEvents(r.Context())

	if err != nil {
		http.Error(w, "failed to query events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "failed to encode events: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
