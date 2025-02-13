package utils

import (
	"github.com/prometheus/client_golang/prometheus"

	"example.com/instrumentation/internal/without/metrics"
)

var (
	// Registry where we are going to register all our metrics for this module
	metricsRegistry *prometheus.Registry

	// All the metrics
	metricsCounter *prometheus.CounterVec
)

func InitMetrics() {
	metricsRegistry = prometheus.NewRegistry()

	metricsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_counter",
			Help: "Count!",
		},
		[]string{"topic"},
	)
	metricsRegistry.MustRegister(metricsCounter)

	metrics.GetManager().RegisterRegistry(metricsRegistry)
}
