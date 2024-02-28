package proxy

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/supporttools/go-sql-proxy/pkg/metrics"
	"github.com/supporttools/go-sql-proxy/pkg/models"
)

var connectionCounter int
var connectionCounterMutex sync.Mutex

// HandleConnection starts the proxy connection, handling data transfer and optional protocol decoding.
func HandleConnection(c *models.Connection) error {
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	mysqlConn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Failed to connect to MySQL [%d]: %s", c.ID, err)
		return err
	}

	// Increment the connection counter
	incrementConnectionCounter()

	defer func() {
		// Decrement the connection counter when the connection is closed
		decrementConnectionCounter()
		mysqlConn.Close()
	}()

	log.Printf("Proxy connection [%d] established to MySQL server at %s", c.ID, address)

	if !c.EnableDecoding {
		return transferData(c, mysqlConn)
	}

	return handleProtocolDecoding(c, mysqlConn)
}

func incrementConnectionCounter() {
	connectionCounterMutex.Lock()
	defer connectionCounterMutex.Unlock()
	connectionCounter++
	metrics.IncrementProxyConnections()
}

func decrementConnectionCounter() {
	connectionCounterMutex.Lock()
	defer connectionCounterMutex.Unlock()
	connectionCounter--
	metrics.DecrementProxyConnections()
}
