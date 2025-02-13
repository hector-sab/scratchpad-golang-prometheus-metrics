package utils

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type CounterCollector struct {
	mu           sync.Mutex
	iteratorDesc *prometheus.Desc
	labelNames   []string
	counts       map[string]int
}

func NewCounterCollector(labelNames []string) *CounterCollector {
	return &CounterCollector{
		iteratorDesc: prometheus.NewDesc(
			"my_counter",
			"Counter!",
			labelNames,
			nil,
		),
		labelNames: labelNames,
		counts:     make(map[string]int),
	}
}

func (c *CounterCollector) Inc(labelValues ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//key := fmt.Sprintf("%v", labelValues)
	key := strings.Join(labelValues, "|")
	slog.Info("Increasing value", "key", key)
	c.counts[key]++
}

func (c *CounterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.iteratorDesc
}

func (c *CounterCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, count := range c.counts {
		// TODO: How do I use fmt.Sscanf...
		//var labelValues []string
		//fmt.Sscanf(key, "[%s]", &labelValues)
		labelValues := strings.Split(key, "|")
		slog.Info("current labels", "key", key, "labels", labelValues, "others", c.labelNames)
		ch <- prometheus.MustNewConstMetric(
			c.iteratorDesc,
			prometheus.CounterValue,
			float64(count),
			labelValues...,
		)
	}
}
