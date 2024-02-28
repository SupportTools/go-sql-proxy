package models

import (
	"context"
)

// Proxy represents the proxy server configuration and state.
type Proxy struct {
	Host           string
	Port           int
	ConnectionID   uint64
	EnableDecoding bool
	Ctx            context.Context
	ShutDownAsked  bool
}
