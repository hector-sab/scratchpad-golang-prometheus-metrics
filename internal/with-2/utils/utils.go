package utils

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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
