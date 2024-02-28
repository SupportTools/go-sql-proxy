package models

import "net"

// Connection represents a proxy connection to a MySQL server.
type Connection struct {
	Host           string
	Port           int
	Conn           net.Conn
	ID             uint64
	EnableDecoding bool
}

func (c *Connection) Read(p []byte) (int, error) {
	return c.Conn.Read(p)
}

func (c *Connection) Write(p []byte) (int, error) {
	return c.Conn.Write(p)
}
