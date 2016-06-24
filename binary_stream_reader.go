package binarystream

import "bytes"
import "fmt"
import "encoding/binary"

type BinaryStream struct {
	buf       []byte
	i         int
	byteOrder binary.ByteOrder
}

func New(buf []byte, byteOrder binary.ByteOrder) *BinaryStream {
	return &BinaryStream{
		buf:       buf,
		byteOrder: byteOrder,
	}
}

func (s *BinaryStream) ReadNullTerminatedString() (string, error) {
	j := bytes.IndexByte(s.buf[s.i:], 0)
	if j < 0 {
		return "", fmt.Errorf("null terminator not found")
	}
	str := string(s.buf[s.i : s.i+j])
	s.i += j + 1
	return str, nil
}

func (s *BinaryStream) ReadUint64() (uint64, error) {
	var val uint64
	buf := bytes.NewReader(s.buf[s.i:])
	err := binary.Read(buf, s.byteOrder, &val)
	if err != nil {
		return 0, err
	}
	s.i += 8
	return val, nil
}

func (s *BinaryStream) Skip(n int) error {
	s.i += n
	if s.i > len(s.buf) {
		return fmt.Errorf("buffer underflow")
	}
	return nil
}

func (s *BinaryStream) ReadFixedString(n int) (string, error) {
	j := s.i + n
	b := s.buf[s.i:j]
	s.i = j
	return string(b), nil
}

func (s *BinaryStream) ReadRemainingString() (string, error) {
	n := len(s.buf) - s.i
	return s.ReadFixedString(n)
}
