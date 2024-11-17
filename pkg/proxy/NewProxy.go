package proxy

import (
	"context"

	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// NewProxy creates a new instance of the Proxy server with optional SSL configuration.
func NewProxy(host string, port int, useSSL bool, ctx context.Context) *models.Proxy {
	return &models.Proxy{
		Host:   host,
		Port:   port,
		UseSSL: useSSL,
		Ctx:    ctx,
	}
}
