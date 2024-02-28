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
		log.Printf("Waiting for shut down signal ^C")
		<-p.Ctx.Done()
		p.ShutDownAsked = true
		log.Printf("Shut down signal received, closing connections...")
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// Error handling
			continue
		}
		p.ConnectionID++
		connection := NewConnection(p.Host, p.Port, conn, p.ConnectionID, p.EnableDecoding)
		go HandleConnection(connection)
	}
}
