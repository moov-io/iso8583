package header

import (
	"fmt"
	"io"
	"strconv"

	"github.com/moov-io/iso8583/encoding"
)

var _ Header = (*BCD2BytesHeader)(nil)

// Visa DEX header is 2 bytes of BCD encoded length
type BCD2BytesHeader struct {
	Len int
}

func NewBCD2BytesHeader() *BCD2BytesHeader {
	return &BCD2BytesHeader{}
}

func (h *BCD2BytesHeader) SetLength(length int) {
	h.Len = length
}

func (h *BCD2BytesHeader) Length() int {
	return h.Len
}

func (h *BCD2BytesHeader) Pack() ([]byte, error) {
	strLen := fmt.Sprintf("%04d", h.Len)
	res, err := encoding.BCD.Encode([]byte(strLen))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *BCD2BytesHeader) Read(reader io.Reader) (int, error) {
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
