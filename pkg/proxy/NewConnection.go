package proxy

import (
	"net"

	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// NewConnection creates a new proxy connection to the target MySQL server.
func NewConnection(host string, port int, conn net.Conn, id uint64, enableDecoding bool) *models.Connection {
	return &models.Connection{
		Host:           host,
		Port:           port,
		Conn:           conn,
		ID:             id,
		EnableDecoding: enableDecoding,
	}
}
