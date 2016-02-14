// Packet with length:
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
	"fmt"
	"io"
	"math"
)

type PacketLenFieldSize uint32

const (
	OneByte   = PacketLenFieldSize(1)
	TwoBytes  = PacketLenFieldSize(2)
	FourBytes = PacketLenFieldSize(4)

	LittleEndian = true
	BigEndian    = false
)

var (
	ErrPacketTooLarge = fmt.Errorf("packet too large")
)

type PacketParser struct {
	packetLenFieldSize PacketLenFieldSize
	packetMaxLen       uint32
	littleEndian       bool
}

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
	OneByteLittleEndian   = NewParser(OneByte, LittleEndian)
	OneByteBigEndian      = NewParser(OneByte, BigEndian)
	TwoBytesLittleEndian  = NewParser(TwoBytes, LittleEndian)
	TwoBytesBigEndian     = NewParser(TwoBytes, BigEndian)
	FourBytesLittleEndian = NewParser(FourBytes, LittleEndian)
	FourBytesBigEndian    = NewParser(FourBytes, BigEndian)
)
