package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"github.com/ccoveille/go-safecast"
)

/*
PacketHeader represents packet header
*/
type PacketHeader struct {
	Length     uint32
	SequenceID uint8
}

/*
InitialHandshakePacket represents initial handshake packet sent by MySQL Server
*/
type InitialHandshakePacket struct {
	ProtocolVersion   uint8
	ServerVersion     []byte
	ConnectionID      uint32
	AuthPluginData    []byte
	Filler            byte
	CapabilitiesFlags CapabilityFlag
	CharacterSet      uint8
	StatusFlags       uint16
	AuthPluginDataLen uint8
	AuthPluginName    []byte
	header            *PacketHeader
}

// Decode decodes the first packet received from the MySQL Server
// It's a handshake packet
func (r *InitialHandshakePacket) Decode(conn net.Conn) error {
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		return err
	}

	header := &PacketHeader{}
	ln := []byte{data[0], data[1], data[2], 0x00}
	header.Length = binary.LittleEndian.Uint32(ln)
	// a single byte integer is the same in BigEndian and LittleEndian
	header.SequenceID = data[3]

	r.header = header
	// Assign payload-only data to new var just for convenience
	payload := data[4 : header.Length+4]
	position := 0

	// Protocol version check
	r.ProtocolVersion = payload[0]
	if r.ProtocolVersion != 0x0a {
		return errors.New("unsupported protocol for the proxy. Only version 10 is supported")
	}
	position++

	// Extract server version
	index := bytes.IndexByte(payload, byte(0x00))
	r.ServerVersion = payload[position:index]
	position = index + 1

	connectionID := payload[position : position+4]
	r.ConnectionID = binary.LittleEndian.Uint32(connectionID)
	position += 4

	// Extract auth plugin data part 1
	r.AuthPluginData = make([]byte, 8)
	copy(r.AuthPluginData, payload[position:position+8])
	position += 8

	r.Filler = payload[position]
	if r.Filler != 0x00 {
		return errors.New("failed to decode filler value")
	}
	position++

	capabilitiesFlags1 := payload[position : position+2]
	position += 2

	r.CharacterSet = payload[position]
	position++

	r.StatusFlags = binary.LittleEndian.Uint16(payload[position : position+2])
	position += 2

	capabilityFlags2 := payload[position : position+2]
	position += 2

	// Reconstruct 32-bit integer from two 16-bit integers
	capLow := binary.LittleEndian.Uint16(capabilitiesFlags1)
	capHi := binary.LittleEndian.Uint16(capabilityFlags2)
	cap := uint32(capLow) | uint32(capHi)<<16
	r.CapabilitiesFlags = CapabilityFlag(cap)

	if r.CapabilitiesFlags&clientPluginAuth != 0 {
		r.AuthPluginDataLen = payload[position]
		if r.AuthPluginDataLen == 0 {
			return errors.New("wrong auth plugin data length")
		}
	}
	position++

	// Skip reserved bytes
	position += 10

	if r.CapabilitiesFlags&clientSecureConn != 0 {
		end := position + Max(13, int(r.AuthPluginDataLen)-8)
		r.AuthPluginData = append(r.AuthPluginData, payload[position:end]...)
		position = end
	}

	index = bytes.IndexByte(payload[position:], byte(0x00))
	if index != -1 {
		r.AuthPluginName = payload[position : position+index]
	} else {
		r.AuthPluginName = payload[position:]
	}

	return nil
}

// Encode encodes the InitialHandshakePacket to a byte slice
func (r InitialHandshakePacket) Encode() ([]byte, error) {
	buf := make([]byte, 0)
	buf = append(buf, r.ProtocolVersion)
	buf = append(buf, r.ServerVersion...)
	buf = append(buf, byte(0x00))

	connectionID := make([]byte, 4)
	binary.LittleEndian.PutUint32(connectionID, r.ConnectionID)
	buf = append(buf, connectionID...)

	auth1 := r.AuthPluginData[0:8]
	buf = append(buf, auth1...)
	buf = append(buf, 0x00)

	cap := make([]byte, 4)
	binary.LittleEndian.PutUint32(cap, uint32(r.CapabilitiesFlags))

	cap1 := cap[0:2]
	cap2 := cap[2:]

	buf = append(buf, cap1...)
	buf = append(buf, r.CharacterSet)

	statusFlag := make([]byte, 2)
	binary.LittleEndian.PutUint16(statusFlag, r.StatusFlags)
	buf = append(buf, statusFlag...)
	buf = append(buf, cap2...)
	buf = append(buf, r.AuthPluginDataLen)

	reserved := make([]byte, 10)
	buf = append(buf, reserved...)
	buf = append(buf, r.AuthPluginData[8:]...)
	buf = append(buf, r.AuthPluginName...)
	buf = append(buf, 0x00)

	// Ensure length fits within uint32 range using safecast
	length, err := safecast.ToUint32(len(buf))
	if err != nil {
		return nil, errors.New("packet size exceeds maximum allowable size for uint32")
	}

	h := PacketHeader{
		Length:     length,
		SequenceID: r.header.SequenceID,
	}

	newBuf := make([]byte, 0, h.Length+4)

	ln := make([]byte, 4)
	binary.LittleEndian.PutUint32(ln, h.Length)

	newBuf = append(newBuf, ln[:3]...)
	newBuf = append(newBuf, h.SequenceID)
	newBuf = append(newBuf, buf...)

	return newBuf, nil
}

// Max returns the larger of x or y.
func (r InitialHandshakePacket) String() string {
	return r.CapabilitiesFlags.String()
}
