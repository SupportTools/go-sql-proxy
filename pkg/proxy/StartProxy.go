package proxy

import (
	"fmt"
	"log"
	"net"

	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// StartProxy starts the proxy server listening for incoming connections.
func StartProxy(p *models.Proxy, port int) error {
	log.Printf("Start listening on: %d", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	go func() {
		log.Printf("Waiting for shutdown signal ^C")
		<-p.Ctx.Done()
		p.ShutDownAsked = true
		log.Printf("Shutdown signal received, closing connections...")
		if err := ln.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			if p.ShutDownAsked {
				return nil // Exit loop if shutdown is requested
			}
			continue
		}
		p.ConnectionID++
		connection := NewConnection(p.Host, p.Port, conn, p.ConnectionID, p.EnableDecoding)

		// Handle connection and log any errors
		go func(c *models.Connection) {
			if err := HandleConnection(c); err != nil {
				log.Printf("Error handling connection %d: %v", c.ID, err)
			}
		}(connection)
	}
}
