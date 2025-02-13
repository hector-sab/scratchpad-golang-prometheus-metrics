package utils

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var counter = NewCounterCollector([]string{"topic"})

// init is a special function in GoLang
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
			fmt.Printf("Iteration: %v\n", i)
			counter.Inc("my-topic")

			time.Sleep(1 * time.Second)
		}
	}
}
