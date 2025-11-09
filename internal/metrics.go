package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Define your metrics here
var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests handled, labeled by route and method",
		},
		[]string{"route", "method"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method"},
	)

	StudentCreationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "student_creations_total",
			Help: "Total number of student records created",
		},
	)

	StatusCodesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_codes_total",
			Help: "Total number of HTTP response status codes",
		},
		[]string{"status_code"},
	)
)
