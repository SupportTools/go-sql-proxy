package proxy

import (
	"context"

	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// NewProxy creates a new instance of the Proxy server.
func NewProxy(ctx context.Context, host string, port int) *models.Proxy {
	return &models.Proxy{
		Host: host,
		Port: port,
		Ctx:  ctx,
	}
}
