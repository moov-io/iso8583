package header

import (
	"fmt"
	"io"
	"strconv"

	"github.com/moov-io/iso8583/encoding"
)

var _ Header = (*VisaDexHeader)(nil)

type VisaDexHeader struct {
	Len int
}

func NewVisaDEXHeader() *VisaDexHeader {
	return &VisaDexHeader{}
}

func (h *VisaDexHeader) SetLength(length int) {
	h.Len = length
}

func (h *VisaDexHeader) Length() int {
	return h.Len
}

func (h *VisaDexHeader) Pack() ([]byte, error) {
	strLen := fmt.Sprintf("%04d", h.Len)
	res, err := encoding.BCD.Encode([]byte(strLen))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *VisaDexHeader) Read(reader io.Reader) (int, error) {
	buf := make([]byte, 2)
	read, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, fmt.Errorf("reading header: %v", err)
	}

	// decode 4 digits from the buf
	bDigits, _, err := encoding.BCD.Decode(buf, 4)
	if err != nil {
		return 0, err
	}

	dataLen, err := strconv.Atoi(string(bDigits))
	if err != nil {
		return 0, fmt.Errorf("converting string to int: %v", err)
	}

	h.Len = dataLen

	return read, nil
}
