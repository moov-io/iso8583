package header

import (
	"fmt"
	"io"
	"strconv"
)

type Header interface {
	Pack() ([]byte, error)
	Read(reader io.Reader) (int, error)
	SetLength(length int)
	Length() int
}

type BaseHeader struct {
	Len int
}

func NewBaseHeader() *BaseHeader {
	return &BaseHeader{}
}

func (h *BaseHeader) SetLength(length int) {
	h.Len = length
}

func (h *BaseHeader) Length() int {
	return h.Len
}

func (h *BaseHeader) Pack() ([]byte, error) {
	return []byte(fmt.Sprintf("%04d", h.Len)), nil
}

func (h *BaseHeader) Read(reader io.Reader) (int, error) {
	buf := make([]byte, 4)
	read, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, fmt.Errorf("reading header: %v", err)
	}

	if read != 4 {
		return 0, fmt.Errorf("excepted to read 4 bytes of the header, got: %v", read)
	}

	l, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, fmt.Errorf("converting header to int: %v", err)
	}
	h.Len = l

	return read, nil
}
