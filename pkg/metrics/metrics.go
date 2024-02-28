package metrics

import (
	"io"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/supporttools/go-sql-proxy/pkg/config"
	"github.com/supporttools/go-sql-proxy/pkg/health"
	"github.com/supporttools/go-sql-proxy/pkg/logging"
)

var logger = logging.SetupLogging()

var (
	// proxy metrics
	proxyConnectionsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "proxy_connections_total",
		Help: "Total number of connections to the proxy.",
	})
	proxyErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "proxy_errors_total",
		Help: "Total number of errors encountered by the proxy.",
	})
	proxyConnectionsOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_connections_open",
		Help: "Number of open connections to the proxy.",
	})

	// data transfer metrics
	DataFromClient = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "proxy_data_from_client_bytes_total",
		Help: "Total number of bytes transferred from client to server through the proxy.",
	})
	DataToClient = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "proxy_data_to_client_bytes_total",
		Help: "Total number of bytes transferred from server to client through the proxy.",
	})

	// proxy request latency
	LastRequestLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_last_request_latency_seconds",
		Help: "The latency of the last proxy request in seconds.",
	})
)

type counterWriter struct {
	counter prometheus.Counter
}

func init() {
	prometheus.MustRegister(DataFromClient)
	prometheus.MustRegister(DataToClient)
}

func StartMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/version", health.VersionHandler())
	mux.HandleFunc("/healthz", health.HealthzHandler(config.CFG.SourceDatabaseUser, config.CFG.SourceDatabasePassword, "localhost", config.CFG.BindPort, config.CFG.SourceDatabaseName))
	mux.HandleFunc("/readyz", health.ReadyzHandler(config.CFG.SourceDatabaseUser, config.CFG.SourceDatabasePassword, config.CFG.SourceDatabaseServer, config.CFG.SourceDatabasePort, config.CFG.SourceDatabaseName))

	serverPortStr := strconv.Itoa(config.CFG.MetricsPort)
	logger.Printf("Metrics server starting on port %d\n", config.CFG.MetricsPort)

	if err := http.ListenAndServe(":"+serverPortStr, mux); err != nil {
		logger.Fatalf("Metrics server failed to start: %v", err)
	}
}

// IncrementProxyConnections increments the proxy connections counter.
func IncrementProxyConnections() {
	proxyConnectionsTotal.Inc()
	proxyConnectionsOpen.Inc()
}

// DecrementProxyConnections decrements the proxy connections counter.
func DecrementProxyConnections() {
	proxyConnectionsOpen.Dec()
}

// IncrementProxyErrors increments the proxy errors counter.
func IncrementProxyErrors() {
	proxyErrors.Inc()
}

// IncrementDataFromClient increments the data from client counter.
func IncrementDataFromClient() {
	DataFromClient.Inc()
}

// IncrementDataToClient increments the data to client counter.
func IncrementDataToClient() {
	DataToClient.Inc()
}

func (cw *counterWriter) Write(p []byte) (int, error) {
	n := len(p)
	cw.counter.Add(float64(n))
	return n, nil
}

// NewCounterWriter creates a new io.Writer that increments the given counter.
func NewCounterWriter(counter prometheus.Counter) io.Writer {
	return &counterWriter{counter}
}
