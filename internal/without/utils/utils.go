package utils

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func Count(iterations int) {
	for i := range iterations {
		fmt.Printf("Iteration: %v\n", i)
		metricsCounter.With(prometheus.Labels{"topic": "my-topic"}).Inc()

		time.Sleep(1 * time.Second)
	}
}
