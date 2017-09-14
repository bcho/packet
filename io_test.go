package packet

import (
	"bytes"
	"math"
	"testing"
)

func testWriter(t *testing.T, packetLenFieldSize PacketLenFieldSize, littleEndian bool, maxPacket uint32) {
	var (
		buf bytes.Buffer
		err error
	)

	writer := NewWriter(&buf, packetLenFieldSize, littleEndian)

	// simple write & read
	data := []byte("hello")
	_, err = writer.Write(data)
	if err != nil {
		t.Error(err)
	}

	// write large packet data
	largeData := make([]byte, uint64(maxPacket)+1)
	_, err = writer.Write(largeData)
	if err != ErrPacketTooLarge {
		t.Errorf("expected error `ErrPacketTooLarge`: %v", err)
	}
}

func TestWriter(t *testing.T) {
	testWriter(t, OneByte, LittleEndian, math.MaxUint8)
	testWriter(t, OneByte, BigEndian, math.MaxUint8)
	testWriter(t, TwoBytes, LittleEndian, math.MaxUint16)
	testWriter(t, TwoBytes, BigEndian, math.MaxUint16)
	testWriter(t, FourBytes, LittleEndian, math.MaxUint32)
	testWriter(t, FourBytes, BigEndian, math.MaxUint32)
}
