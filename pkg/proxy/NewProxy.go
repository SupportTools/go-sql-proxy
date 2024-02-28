package proxy

import (
	"context"

	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// NewProxy creates a new instance of the Proxy server.
func NewProxy(host string, port int, ctx context.Context) *models.Proxy {
	return &models.Proxy{
		Host: host,
		Port: port,
		Ctx:  ctx,
	}
}
