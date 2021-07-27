package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/moov-io/iso8583/encoding"
)

const (
	sessionControlIndicator = byte('2')
	MaxMessageLength        = 2048
)

type VMLH struct {
	Len uint16
	//
	IsSessionControl bool
}

func NewVMLHeader() *VMLH {
	return &VMLH{}
}

func (h *VMLH) SetLength(length int) error {
	if length > math.MaxUint16 {
		return fmt.Errorf("length %d exceeds max length for 2 bytes header %d", length, math.MaxUint16)
	}

	h.Len = uint16(length)

	return nil
}

func (h *VMLH) Length() int {
	return int(h.Len)
}

func (h *VMLH) WriteTo(w io.Writer) (int, error) {
	err := binary.Write(w, binary.BigEndian, h.Len)
	if err != nil {
		return 0, fmt.Errorf("wrigint uint16 into writer: %v", err)
	}

	return binary.Size(h.Len), nil
}

func (h *VMLH) ReadFrom(r io.Reader) (int, error) {
	header := make([]byte, 4)

	// read full header
	read, err := io.ReadFull(r, header)
	if err != nil {
		return 0, fmt.Errorf("reading 4 bytes from reader: %v", err)
	}

	// read 2 bytes length
	err = binary.Read(bytes.NewReader(header), binary.BigEndian, &h.Len)
	if err != nil {
		return 0, fmt.Errorf("reading uint16 length from reader: %v", err)
	}

	// read message format and platform
	indicators, _, err := encoding.BCD.Decode(header[3:], 2)
	if err != nil {
		return 0, fmt.Errorf("decoding indicators: %v", err)
	}

	h.IsSessionControl = (indicators[0] == sessionControlIndicator)

	return read, nil
}
