package iso8583

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBCDDecode(t *testing.T) {

	b := []byte("954")
	r := rbcd(b)
	assert.Equal(t, "0954", fmt.Sprintf("%X", r))

	r = lbcd(b)
	assert.Equal(t, "9540", fmt.Sprintf("%X", r))

	b = []byte("31")
	r = lbcd(b)
	assert.Equal(t, "31", fmt.Sprintf("%X", r))
	r = rbcd(b)
	assert.Equal(t, "31", fmt.Sprintf("%X", r))

	b = []byte("123ab4")
	assert.Equal(t, []byte("\x12\x3a\xb4"), bcd(b))
	b = []byte("00")
	assert.Equal(t, []byte("\x00"), bcd(b))

	assert.Panics(t,
		func() {
			bcd([]byte("test"))
		}, "Calling bcd() with invalid hex should panic")

}

func TestBCDEncode(t *testing.T) {
	assert.Equal(t, []byte("12a34f"), bcd2Ascii([]byte("\x12\xa3\x4f")))

	assert.Equal(t, []byte("12345"), bcdl2Ascii([]byte("\x12\x34\x50"), 5))

	assert.Equal(t, []byte("12345"), bcdr2Ascii([]byte("\x01\x23\x45"), 5))
}
