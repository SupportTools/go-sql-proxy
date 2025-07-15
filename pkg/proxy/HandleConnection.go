package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"log"
	"net"

	"github.com/supporttools/go-sql-proxy/pkg/config"
	"github.com/supporttools/go-sql-proxy/pkg/metrics"
	"github.com/supporttools/go-sql-proxy/pkg/models"
)

// HandleConnection starts the proxy connection, handling data transfer and optional protocol decoding.
func HandleConnection(c *models.Connection) error {
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	
	var mysqlConn net.Conn
	var err error
	
	if config.CFG.UseSSL {
		mysqlConn, err = dialWithSSL(address)
	} else {
		mysqlConn, err = net.Dial("tcp", address)
	}
	
	if err != nil {
		log.Printf("Failed to connect to MySQL [%d]: %s", c.ID, err)
		return err
	}

	metrics.IncrementProxyConnections() // Increment metric counter

	defer func() {
		metrics.DecrementProxyConnections() // Decrement metric counter when connection is closed
		if err := mysqlConn.Close(); err != nil {
			log.Printf("Error closing MySQL connection [%d]: %v", c.ID, err)
		}
	}()

	log.Printf("Proxy connection [%d] established to MySQL server at %s", c.ID, address)

	if !c.EnableDecoding {
		return transferData(c, mysqlConn)
	}

	return handleProtocolDecoding(c, mysqlConn)
}

// dialWithSSL creates an SSL/TLS connection to the MySQL server
func dialWithSSL(address string) (net.Conn, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.CFG.SSLSkipVerify, // #nosec G402 - InsecureSkipVerify is configurable for development environments
	}

	// Load custom CA if provided
	if config.CFG.SSLCAFile != "" {
		caCert, err := os.ReadFile(config.CFG.SSLCAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	// Load client certificates if provided
	if config.CFG.SSLCertFile != "" && config.CFG.SSLKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(config.CFG.SSLCertFile, config.CFG.SSLKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificates: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tls.Dial("tcp", address, tlsConfig)
}
