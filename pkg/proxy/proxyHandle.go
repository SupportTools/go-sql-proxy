package proxy

import (
	"log"
	"net"

	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// proxyHandle creates and manages a new connection.
func proxyHandle(p *models.Proxy, conn net.Conn) {
	connection := NewConnection(p.Host, p.Port, conn, p.ConnectionID, p.EnableDecoding)
	err := HandleConnection(connection)
	if err != nil {
		log.Printf("Error handling proxy connection: %s", err)
	}
}
