package handler

import (
	"github.com/diemensa/event-analytics-service/internal/metrics"
	"github.com/diemensa/event-analytics-service/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func SetupHandlers(s service.Service) {
	eventHandler := NewEventHandler(s)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			eventHandler.HandleGetEvents(w, r)
		case http.MethodPost:
			eventHandler.HandleCreateEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	metrics.PrometheusInit()
	http.Handle("/metrics", promhttp.Handler())
}
