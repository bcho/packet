package packet

import (
	"bytes"
	"math"
	"testing"
)

func testParser(t *testing.T, parser PacketParser, maxPacket uint32) {
	var (
		buf bytes.Buffer
		err error
	)

	// simple write & read
	data := []byte("hello")
	err = parser.Write(&buf, data)
	if err != nil {
		t.Error(err)
	}

	readData, err := parser.Read(&buf)
	if err != nil {
		t.Error(err)
	}
	if string(readData) != string(data) {
		t.Errorf("unexpected data read: %s", string(readData))
	}

	// write large packet data
	largeData := make([]byte, uint64(maxPacket)+1)
	err = parser.Write(&buf, largeData)
	if err != ErrPacketTooLarge {
		t.Errorf("expected error `ErrPacketTooLarge`: %v", err)
	}
}

func TestOneByte(t *testing.T) {
	testParser(t, OneByteLittleEndian, math.MaxUint8)
	testParser(t, OneByteBigEndian, math.MaxUint8)
}

func TestTwoBytes(t *testing.T) {
	testParser(t, TwoBytesLittleEndian, math.MaxUint16)
	testParser(t, TwoBytesBigEndian, math.MaxUint16)
}

func TestFourBytes(t *testing.T) {
	testParser(t, FourBytesLittleEndian, math.MaxUint32)
	testParser(t, FourBytesBigEndian, math.MaxUint32)
}

func TestReadLargePacket(t *testing.T) {
	parser := NewParser(OneByte, BigEndian)
	parser.packetMaxLen = 10 // make error happen

	largeData := make([]byte, 1)
	largeData[0] = byte(100)
	_, err := parser.Read(bytes.NewBuffer(largeData))
	if err != ErrPacketTooLarge {
		t.Errorf("expected error `ErrPacketTooLarge`: %v", err)
	}
}

func benchmarkParser(b *testing.B, parser PacketParser) {
	var buf bytes.Buffer
	data := []byte("hello")

	for i := 0; i < b.N; i++ {
		parser.Write(&buf, data)
		parser.Read(&buf)
	}
}

func BenchmarkOneByteLittleEndian(b *testing.B) {
	benchmarkParser(b, OneByteLittleEndian)
}

func BenchmarkOneByteBigEndian(b *testing.B) {
	benchmarkParser(b, OneByteBigEndian)
}

func BenchmarkTwoBytesLittleEndian(b *testing.B) {
	benchmarkParser(b, TwoBytesLittleEndian)
}

func BenchmarkTwoBytesBigEndian(b *testing.B) {
	benchmarkParser(b, TwoBytesBigEndian)
}

func BenchmarkFourBytesLittleEndian(b *testing.B) {
	benchmarkParser(b, FourBytesLittleEndian)
}

func BenchmarkFourBytesBigEndian(b *testing.B) {
	benchmarkParser(b, FourBytesBigEndian)
}
