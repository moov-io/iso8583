package header

import "fmt"

type Header interface {
	Pack() ([]byte, error)
	Unpack(data []byte) (int, error)
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

func (h *BaseHeader) Unpack(data []byte) (int, error) {
	return 0, nil
}
