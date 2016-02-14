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
	"encoding/binary"
	"errors"
	"io"
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

// A PacketParser is used for parsing or building a packet.
type PacketParser struct {
	packetLenFieldSize PacketLenFieldSize
	packetMaxLen       uint32
	littleEndian       bool
}

// Instance a new packet parser.
func NewParser(packetLenFieldSize PacketLenFieldSize, littleEndian bool) PacketParser {
	p := PacketParser{
		packetLenFieldSize: packetLenFieldSize,
		littleEndian:       littleEndian,
	}

	switch packetLenFieldSize {
	case OneByte:
		p.packetMaxLen = math.MaxUint8
	case TwoBytes:
		p.packetMaxLen = math.MaxUint16
	case FourBytes:
		p.packetMaxLen = math.MaxUint32
	default:
		p.packetMaxLen = math.MaxUint32
	}

	return p
}

// Read reads a packet from reader. Any error encountered during the read will be returned.
// If packet is too large, ErrPacketTooLarge will be returned.
func (p PacketParser) Read(reader io.Reader) ([]byte, error) {
	packetLenBuffer := make([]byte, p.packetLenFieldSize)
	if _, err := io.ReadFull(reader, packetLenBuffer); err != nil {
		return nil, err
	}

	var packetLen uint32
	switch p.packetLenFieldSize {
	case OneByte:
		packetLen = uint32(packetLenBuffer[0])
	case TwoBytes:
		if p.littleEndian {
			packetLen = uint32(binary.LittleEndian.Uint16(packetLenBuffer))
		} else {
			packetLen = uint32(binary.BigEndian.Uint16(packetLenBuffer))
		}
	case FourBytes:
		if p.littleEndian {
			packetLen = binary.LittleEndian.Uint32(packetLenBuffer)
		} else {
			packetLen = binary.BigEndian.Uint32(packetLenBuffer)
		}
	}

	if packetLen > p.packetMaxLen {
		return nil, ErrPacketTooLarge
	}

	packetData := make([]byte, packetLen)
	if _, err := io.ReadFull(reader, packetData); err != nil {
		return nil, err
	}

	return packetData, nil
}

// Write writes a packet data with length to writer. Any error encountered during the write will be returned.
// If packet data is too large, ErrPacketTooLarge will be returned.
func (p PacketParser) Write(writer io.Writer, data []byte) error {
	packetLen := uint32(len(data))
	overflow := len(data) != int(packetLen) // FIXME use more robust mean
	if packetLen > p.packetMaxLen || overflow {
		return ErrPacketTooLarge
	}

	packet := make([]byte, uint32(p.packetLenFieldSize)+packetLen)
	switch p.packetLenFieldSize {
	case OneByte:
		packet[0] = byte(packetLen)
	case TwoBytes:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(packet, uint16(packetLen))
		} else {
			binary.BigEndian.PutUint16(packet, uint16(packetLen))
		}
	case FourBytes:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(packet, packetLen)
		} else {
			binary.BigEndian.PutUint32(packet, packetLen)
		}
	}
	copy(packet[p.packetLenFieldSize:], data)

	_, err := writer.Write(packet)
	return err
}

var (
	OneByteLittleEndian   = NewParser(OneByte, LittleEndian)   // 8 bits length (0 ~ 255), little endian
	OneByteBigEndian      = NewParser(OneByte, BigEndian)      // 8 bits length (0 ~ 255), big endian
	TwoBytesLittleEndian  = NewParser(TwoBytes, LittleEndian)  // 16 bits length (0 ~ 65535), little endian
	TwoBytesBigEndian     = NewParser(TwoBytes, BigEndian)     // 16 bits length (0 ~ 65535), big endian
	FourBytesLittleEndian = NewParser(FourBytes, LittleEndian) // 32 bits length (0 ~ 4294967295), little endian
	FourBytesBigEndian    = NewParser(FourBytes, BigEndian)    // 32 bits length (0 ~ 4294967295), big endian
)
