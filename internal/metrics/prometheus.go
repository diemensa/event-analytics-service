package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MessagesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_processed_total",
			Help: "Total number of processed events",
		},
		[]string{"status"},
	)

	MessageProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "event_processing_duration",
			Help:    "Duration of event processing (in seconds)",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func PrometheusInit() {
	prometheus.MustRegister(MessagesProcessed)
	prometheus.MustRegister(MessageProcessingDuration)
}
