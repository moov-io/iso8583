package network

import (
	"fmt"
	"io"
	"strconv"
)

var _ Header = (*ASCII4BytesHeader)(nil)

// ASCII4BytesHeader is 4 bytes ASCII encoded length
type ASCII4BytesHeader struct {
	Len int
}

func NewASCII4BytesHeader() *ASCII4BytesHeader {
	return &ASCII4BytesHeader{}
}

func (h *ASCII4BytesHeader) SetLength(length int) {
	h.Len = length
}

func (h *ASCII4BytesHeader) Length() int {
	return h.Len
}

func (h *ASCII4BytesHeader) WriteTo(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "%04d", h.Len)
}

func (h *ASCII4BytesHeader) ReadFrom(r io.Reader) (int, error) {
	buf := make([]byte, 4)
	read, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, fmt.Errorf("reading header: %v", err)
	}

	if read != 4 {
		return 0, fmt.Errorf("expected to read 4 bytes of the header, got: %v", read)
	}

	l, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, fmt.Errorf("converting header to int: %v", err)
	}
	h.Len = l

	return read, nil
}
