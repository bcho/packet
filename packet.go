// Packet packet helps create packet with length.
//
// Packet structure:
//
//      --------------
//      | len | data |
//      --------------
//
// Read a packet:
//
//      data, _ := FourBytesBigEndian.Read(tcpConn)
//      // process data
//
// Write a packet:
//
//      FourBytesBigEndian.Write(tcpConn, []byte("hello, world"))
package packet

import (
	"errors"
	"math"
)

// Packet length field size.
type PacketLenFieldSize uint32

const (
	OneByte   = PacketLenFieldSize(1) // use 8 bites for length field
	TwoBytes  = PacketLenFieldSize(2) // use 16 bits for length field
	FourBytes = PacketLenFieldSize(4) // use 32 bits for length field

	LittleEndian = true  // use little endian
	BigEndian    = false // use big endian
)

var (
	// ErrPacketTooLarge will be used when packet data too large.
	ErrPacketTooLarge = errors.New("packet too large")
)

func maxLenOfPacket(packetLenFieldSize PacketLenFieldSize) uint32 {
	switch packetLenFieldSize {
	case OneByte:
		return math.MaxUint8
	case TwoBytes:
		return math.MaxUint16
	case FourBytes:
		return math.MaxUint32
	default:
		return math.MaxUint32
	}
}
