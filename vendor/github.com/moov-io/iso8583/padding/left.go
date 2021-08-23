package padding

import (
	"bytes"
	"unicode/utf8"
)

var Left func(pad rune) Padder = NewLeftPadder

type leftPadder struct {
	pad []byte
}

func NewLeftPadder(pad rune) Padder {
	buf := make([]byte, utf8.RuneLen(pad))
	utf8.EncodeRune(buf, pad)

	return &leftPadder{buf}
}

func (p *leftPadder) Pad(data []byte, length int) []byte {
	if len(data) >= length {
		return data
	}

	padding := bytes.Repeat(p.pad, length-len(data))
	return append(padding, data...)
}

func (p *leftPadder) Unpad(data []byte) []byte {
	pad, _ := utf8.DecodeRune(p.pad)

	return bytes.TrimLeftFunc(data, func(r rune) bool {
		return r == pad
	})
}

func (p *leftPadder) Inspect() []byte {
	return p.pad
}
