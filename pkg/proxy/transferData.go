package proxy

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/supporttools/go-sql-proxy/pkg/metrics"
	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// transferData handles the bi-directional transfer of data between the client and the server.
// It also measures and updates the latency of the last request using a Prometheus gauge.
func transferData(c *models.Connection, conn net.Conn) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Channel to capture latencies from both directions.
	latencies := make(chan float64, 2)

	// Transfer data from client to server
	go func() {
		defer wg.Done()
		startTime := time.Now()
		if _, err := io.Copy(conn, io.TeeReader(c.Conn, metrics.NewCounterWriter(metrics.DataFromClient))); err != nil {
			log.Printf("Error transferring data from client to server: %v", err)
		}
		latency := time.Since(startTime).Seconds()
		latencies <- latency // Send latency to channel
	}()

	// Transfer data from server to client
	go func() {
		defer wg.Done()
		startTime := time.Now()
		if _, err := io.Copy(c.Conn, io.TeeReader(conn, metrics.NewCounterWriter(metrics.DataToClient))); err != nil {
			log.Printf("Error transferring data from server to client: %v", err)
		}
		latency := time.Since(startTime).Seconds()
		latencies <- latency // Send latency to channel
	}()

	// Wait for both transfers to complete
	wg.Wait()
	close(latencies) // Close the channel

	// Calculate the max latency of the two directions as the last request latency
	var maxLatency float64
	for latency := range latencies {
		if latency > maxLatency {
			maxLatency = latency
		}
	}

	// Update the gauge with the latency of the last request
	metrics.LastRequestLatency.Set(maxLatency)

	return nil
}
