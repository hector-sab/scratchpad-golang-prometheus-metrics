package utils

import (
	"sync"

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
