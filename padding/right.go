package padding

import (
	"bytes"
	"unicode/utf8"
)

// Right returns a new right-side padder
var Right func(pad rune) Padder = NewRightPadder

type rightPadder struct {
	pad []byte
}

// NewRightPadder takes the given byte character and returns a padder
// which pads fields to the right of their values (for left-justified values)
func NewRightPadder(pad rune) Padder {
	buf := make([]byte, utf8.RuneLen(pad))
	utf8.EncodeRune(buf, pad)

	return &rightPadder{buf}
}

func (p *rightPadder) Pad(data []byte, length int) []byte {
	if len(data) >= length {
		return data
	}

	padding := bytes.Repeat(p.pad, length-len(data))
	return append(data, padding...)
}

func (p *rightPadder) Unpad(data []byte) []byte {
	pad, _ := utf8.DecodeRune(p.pad)

	return bytes.TrimRightFunc(data, func(r rune) bool {
		return r == pad
	})
}

func (p *rightPadder) Inspect() []byte {
	return p.pad
}
