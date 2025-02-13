package metrics

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	instance        *Manager
	once            sync.Once
	metricsRegistry *prometheus.Registry
)

type Manager struct {
	registries []prometheus.Gatherer
	mu         sync.RWMutex
}

// CreateAndStartServer creates and HTTP server and initializes it in the background. Make sure you have
// registered all the prometheus registries with Manager.RegisterRegistry before creating this server.
// Otherwise the metrics won't be available though the adderss specified.
func CreateAndStartServer(address string, port int, endpoint string) *http.Server {
	http.Handle(endpoint, GetManager().Handler())
	addr := fmt.Sprintf("%v:%v", address, port)
	svr := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to initialize HTTP server for metrics",
				"error", err,
				"address", addr,
				"endpoint", endpoint,
			)
			os.Exit(1)
		}
	}()
	return svr
}

func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			registries: make([]prometheus.Gatherer, 0),
		}
	})
	return instance
}

// RegisterRegistry registers new metrics registries. Make sure register any new registry before
// the HTTP is initialized. Otherwise the metrics won't be available.
func (m *Manager) RegisterRegistry(registry *prometheus.Registry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registries = append(m.registries, registry)
}

func (m *Manager) Handler() http.Handler {
	return promhttp.HandlerFor(
		prometheus.Gatherers(m.registries),
		promhttp.HandlerOpts{},
	)
}

// InitMetrics initializes default go metrics of the process.
func InitMetrics() {
	metricsRegistry = prometheus.NewRegistry()

	metricsRegistry.MustRegister(collectors.NewGoCollector())
	metricsRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	GetManager().RegisterRegistry(metricsRegistry)
}
