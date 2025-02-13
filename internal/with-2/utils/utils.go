package utils

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type CounterCollector struct {
	counter *prometheus.CounterVec
	mu      sync.Mutex
}

func NewCounterCollector() *CounterCollector {
	return &CounterCollector{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "my_counter",
				Help: "Counter!",
			},
			[]string{"topic"},
		),
	}
}

func (c *CounterCollector) Describe(ch chan<- *prometheus.Desc) {
	c.counter.Describe(ch)
}

func (c *CounterCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.counter.Collect(ch)
}

var counter = NewCounterCollector()

func init() {
	prometheus.MustRegister(counter)
}

func Count(ctx context.Context, iterations int) {
	for i := range iterations {
		select {
		case <-ctx.Done():
			slog.Debug("Stopping count. Context done")
			return
		default:
			slog.Info("Counting!", "iteration", i)
			counter.counter.With(prometheus.Labels{"topic": "my-topic"}).Inc()

			time.Sleep(1 * time.Second)
		}
	}
}
