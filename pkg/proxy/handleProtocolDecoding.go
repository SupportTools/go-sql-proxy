package proxy

import (
	"log"
	"net"

	"github.com/supporttools/go-sql-proxy/pkg/models"
	"github.com/supporttools/go-sql-proxy/pkg/protocol"
)

// handleProtocolDecoding decodes the MySQL protocol if enabled, starting with the handshake packet.
func handleProtocolDecoding(c *models.Connection, mysqlConn net.Conn) error {
	handshakePacket := &protocol.InitialHandshakePacket{}
	if err := handshakePacket.Decode(mysqlConn); err != nil {
		log.Printf("Failed to decode handshake packet [%d]: %s", c.ID, err)
		return err
	}

	//log.Printf("Decoded InitialHandshakePacket for connection [%d]: %+v", c.ID, handshakePacket)

	response, err := handshakePacket.Encode()
	if err != nil {
		log.Printf("Failed to encode handshake response [%d]: %s", c.ID, err)
		return err
	}

	if _, err := c.Conn.Write(response); err != nil {
		log.Printf("Failed to send handshake response to client [%d]: %s", c.ID, err)
		return err
	}

	return transferData(c, mysqlConn)
}
