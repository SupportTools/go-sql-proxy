package proxy

import "github.com/supporttools/go-sql-proxy/pkg/models"

// EnableDecoding enables protocol decoding for the proxy.
func EnableDecoding(p *models.Proxy) {
	p.EnableDecoding = true
}
