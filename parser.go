package packet

import (
	"io"
	"math"
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
	return readPacket(
		reader,
		p.littleEndian,
		p.packetLenFieldSize,
		p.packetMaxLen,
	)
}

// Write writes a packet data with length to writer. Any error encountered during the write will be returned.
// If packet data is too large, ErrPacketTooLarge will be returned.
func (p PacketParser) Write(writer io.Writer, data []byte) error {
	_, err := writePacket(
		writer,
		data,
		p.littleEndian,
		p.packetLenFieldSize,
		p.packetMaxLen,
	)

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
