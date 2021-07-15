package network

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

var _ Header = (*Hex4BytesHeader)(nil)

// Hex4BytesHeader is 4 bytes HEX encoded length
type Hex4BytesHeader struct {
	Len int
}

func NewHex4BytesHeader() *Hex4BytesHeader {
	return &Hex4BytesHeader{}
}

func (h *Hex4BytesHeader) SetLength(length int) {
	h.Len = length
}

func (h *Hex4BytesHeader) Length() int {
	return h.Len
}

func (h *Hex4BytesHeader) WriteTo(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "%04X", h.Len)
}

func (h *Hex4BytesHeader) ReadFrom(r io.Reader) (int, error) {
	buf := make([]byte, 4)
	read, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, fmt.Errorf("reading header: %v", err)
	}

	if read != 4 {
		return 0, fmt.Errorf("expected to read 4 bytes of the header, got: %v", read)
	}

	parsed, err := strconv.ParseInt(string(buf), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("converting hex to int: %v", err)
	}

	if parsed < 0 || parsed > math.MaxInt32 {
		return 0, fmt.Errorf("converting parsed integer into smaller bit size than expected: %d", parsed)
	}

	h.Len = int(parsed)

	return read, nil
}
