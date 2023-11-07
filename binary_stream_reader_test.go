package binarystream_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/alexcb/binarystream"
)

func TestReadNullTerminatedString(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.LittleEndian)

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "")

	s, err = bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "hello")
}

func TestReadUint8(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.LittleEndian)

	x, err := bsr.ReadUint8()
	NoError(t, err)
	Equal(t, x, uint8(0))

	x, err = bsr.ReadUint8()
	NoError(t, err)
	Equal(t, x, uint8(104))

	x, err = bsr.ReadUint8()
	NoError(t, err)
	Equal(t, x, uint8(101))

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "llo")
}

func TestReadUint16LittleEndian(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.LittleEndian)

	x, err := bsr.ReadUint16()
	NoError(t, err)
	Equal(t, x, uint16(26624))

	x, err = bsr.ReadUint16()
	NoError(t, err)
	Equal(t, x, uint16(27749))

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "lo")
}

func TestReadUint16BigEndian(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.BigEndian)

	x, err := bsr.ReadUint16()
	NoError(t, err)
	Equal(t, x, uint16(104))

	x, err = bsr.ReadUint16()
	NoError(t, err)
	Equal(t, x, uint16(25964))

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "lo")
}

func TestReadUint32LittleEndian(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.LittleEndian)

	x, err := bsr.ReadUint32()
	NoError(t, err)
	Equal(t, x, uint32(1818585088))

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "lo")
}

func TestReadUint64LittleEndianUnderflow(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.LittleEndian)

	_, err := bsr.ReadUint64()
	Error(t, err, binarystream.ErrBufferUnderflow)
}

func TestReadUint64LittleEndian(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.LittleEndian)

	x, err := bsr.ReadUint64()
	NoError(t, err, nil)
	Equal(t, x, uint64(0x6C6C656800000000))

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "o")
}

func TestReadUint64BigEndian(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.BigEndian)

	x, err := bsr.ReadUint64()
	NoError(t, err, nil)
	Equal(t, x, uint64(1751477356))

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "o")
}

func TestReadBytes(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x01, 0x00, 0x00, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00}, binary.BigEndian)

	b, err := bsr.ReadBytes(5)
	NoError(t, err, nil)
	Equal(t, b, []byte{0x00, 0x01, 0x00, 0x00, 0x68})

	s, err := bsr.ReadNullTerminatedString()
	NoError(t, err)
	Equal(t, s, "ello")

	// make sure it didn't change
	Equal(t, b, []byte{0x00, 0x01, 0x00, 0x00, 0x68})
}

func TestReadPrefixedString(t *testing.T) {
	bsr := binarystream.NewReaderFromBytes([]byte{0x00, 0x01, 0x00, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x01, 0x21, 0x01}, binary.BigEndian)

	s, err := bsr.ReadUint8PrefixedString()
	NoError(t, err, nil)
	Equal(t, s, "")

	s, err = bsr.ReadUint8PrefixedString()
	NoError(t, err, nil)
	Equal(t, s, string([]byte{0x00}))

	s, err = bsr.ReadUint8PrefixedString()
	NoError(t, err, nil)
	Equal(t, s, "hello")

	s, err = bsr.ReadUint8PrefixedString()
	NoError(t, err, nil)
	Equal(t, s, "!")

	_, err = bsr.ReadUint8PrefixedString()
	Error(t, err, binarystream.ErrBufferUnderflow)
}

func TestUnderrunIsRecoverable(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	bsr := binarystream.NewReader(buf, binary.LittleEndian)

	_, err := bsr.ReadUint16PrefixedString()
	Error(t, err, binarystream.ErrBufferUnderflow)

	_, err = buf.Write([]byte{0x03, 0x00}) // 3 chars
	NoError(t, err)

	_, err = bsr.ReadUint16PrefixedString()
	Error(t, err, binarystream.ErrBufferUnderflow)

	_, err = buf.Write([]byte{0x61}) // a
	NoError(t, err)

	_, err = bsr.ReadUint16PrefixedString()
	Error(t, err, binarystream.ErrBufferUnderflow)

	_, err = buf.Write([]byte{0x62, 0x63}) // b c
	NoError(t, err)

	s, err := bsr.ReadUint16PrefixedString()
	NoError(t, err)
	Equal(t, s, "abc")
}
