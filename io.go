package packet

import (
	"encoding/binary"
	"io"
)

func readPacket(
	reader io.Reader,
	littleEndian bool,
	packetLenFieldSize PacketLenFieldSize,
	packetMaxLen uint32,
) ([]byte, error) {
	packetLenBuffer := make([]byte, packetLenFieldSize)
	if _, err := io.ReadFull(reader, packetLenBuffer); err != nil {
		return nil, err
	}

	var packetLen uint32
	switch packetLenFieldSize {
	case OneByte:
		packetLen = uint32(packetLenBuffer[0])
	case TwoBytes:
		if littleEndian {
			packetLen = uint32(binary.LittleEndian.Uint16(packetLenBuffer))
		} else {
			packetLen = uint32(binary.BigEndian.Uint16(packetLenBuffer))
		}
	case FourBytes:
		if littleEndian {
			packetLen = binary.LittleEndian.Uint32(packetLenBuffer)
		} else {
			packetLen = binary.BigEndian.Uint32(packetLenBuffer)
		}
	}

	if packetLen > packetMaxLen {
		return nil, ErrPacketTooLarge
	}

	packetData := make([]byte, packetLen)
	if _, err := io.ReadFull(reader, packetData); err != nil {
		return nil, err
	}

	return packetData, nil
}

// ReadPacket reads a packet from reader. Any error encountered during the read will be returned.
// If packet is too large, ErrPacketTooLarge will be returned.
func ReadPacket(reader io.Reader, packetLenFieldSize PacketLenFieldSize, littleEndian bool) ([]byte, error) {
	return readPacket(
		reader,
		littleEndian,
		packetLenFieldSize,
		maxLenOfPacket(packetLenFieldSize),
	)
}

func writePacket(
	writer io.Writer,
	data []byte,
	littleEndian bool,
	packetLenFieldSize PacketLenFieldSize,
	packetMaxLen uint32,
) (int, error) {
	packetLen := uint32(len(data))
	overflow := len(data) != int(packetLen) // FIXME use more robust mean

	if packetLen > packetMaxLen || overflow {
		return 0, ErrPacketTooLarge
	}

	packet := make([]byte, uint32(packetLenFieldSize)+packetLen)
	switch packetLenFieldSize {
	case OneByte:
		packet[0] = byte(packetLen)
	case TwoBytes:
		if littleEndian {
			binary.LittleEndian.PutUint16(packet, uint16(packetLen))
		} else {
			binary.BigEndian.PutUint16(packet, uint16(packetLen))
		}
	case FourBytes:
		if littleEndian {
			binary.LittleEndian.PutUint32(packet, packetLen)
		} else {
			binary.BigEndian.PutUint32(packet, packetLen)
		}
	}
	copy(packet[packetLenFieldSize:], data)

	return writer.Write(packet)
}

// WritePacket writes a packet to writer with given settings. Any error encountered during the write will be returned.
// If packet data is too large, ErrPacketTooLarge will be returned.
func WritePacket(writer io.Writer, data []byte, littleEndian bool, packetLenFieldSize PacketLenFieldSize) (int, error) {
	return writePacket(
		writer,
		data,
		littleEndian,
		packetLenFieldSize,
		maxLenOfPacket(packetLenFieldSize),
	)
}
