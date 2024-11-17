package proxy

import (
	"fmt"
	"log"
	"net"

	"github.com/supporttools/go-sql-proxy/pkg/metrics"
	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// HandleConnection starts the proxy connection, handling data transfer and optional protocol decoding.
func HandleConnection(c *models.Connection) error {
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	mysqlConn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Failed to connect to MySQL [%d]: %s", c.ID, err)
		return err
	}

	metrics.IncrementProxyConnections() // Increment metric counter

	defer func() {
		metrics.DecrementProxyConnections() // Decrement metric counter when connection is closed
		if err := mysqlConn.Close(); err != nil {
			log.Printf("Error closing MySQL connection [%d]: %v", c.ID, err)
		}
	}()

	log.Printf("Proxy connection [%d] established to MySQL server at %s", c.ID, address)

	if !c.EnableDecoding {
		return transferData(c, mysqlConn)
	}

	return handleProtocolDecoding(c, mysqlConn)
}
