package handler

import (
	"github.com/diemensa/event-analytics-service/internal/metrics"
	"github.com/diemensa/event-analytics-service/internal/service"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func SetupHandlers(s service.Service) *mux.Router {
	eventHandler := NewEventHandler(s)
	metrics.PrometheusInit()

	r := mux.NewRouter()

	r.HandleFunc("/events", eventHandler.HandleCreateEvent).Methods(http.MethodPost)
	r.HandleFunc("/events", eventHandler.HandleGetEvents).Methods(http.MethodGet)
	r.Handle("/metrics", promhttp.Handler())

	return r
}
