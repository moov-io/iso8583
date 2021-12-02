package network

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Binary2Bytes struct {
	Len uint16
}

func NewBinary2BytesHeader() *Binary2Bytes {
	return &Binary2Bytes{}
}

func (h *Binary2Bytes) SetLength(length int) error {
	if length > math.MaxUint16 {
		return fmt.Errorf("length %d exceeds max length for 2 bytes header %d", length, math.MaxUint16)
	}

	h.Len = uint16(length)

	return nil
}

func (h *Binary2Bytes) Length() int {
	return int(h.Len)
}

func (h *Binary2Bytes) WriteTo(w io.Writer) (int, error) {
	err := binary.Write(w, binary.BigEndian, h.Len)
	if err != nil {
		return 0, fmt.Errorf("wrigint uint16 into writer: %w", err)
	}

	return binary.Size(h.Len), nil
}

func (h *Binary2Bytes) ReadFrom(r io.Reader) (int, error) {
	err := binary.Read(r, binary.BigEndian, &h.Len)
	if err != nil {
		return 0, fmt.Errorf("reading uint16 from reader: %w", err)
	}

	return binary.Size(h.Len), nil
}
