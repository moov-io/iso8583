package utils

// Partially copied from https://github.com/boljen/go-bitmap (MIT license)

import (
	"fmt"
)

type Bitmap struct {
	bitmap []byte
}

var (
	tA = [8]byte{128, 64, 32, 16, 8, 4, 2, 1}
)

func NewBitmap(l int) *Bitmap {
	remainder := l % 8
	if remainder != 0 {
		remainder = 1
	}

	return &Bitmap{
		bitmap: make([]byte, l/8+remainder),
	}
}

func NewBitmapFromData(data []byte) *Bitmap {
	return &Bitmap{
		bitmap: data,
	}
}

func (b *Bitmap) IsSet(i int) bool {
	i = i - 1
	return b.bitmap[i/8]&tA[i%8] != 0
}

func (b *Bitmap) Set(i int) {
	i = i - 1
	index := i / 8
	bit := i % 8
	b.bitmap[index] = b.bitmap[index] | tA[bit]
}

func (b *Bitmap) Bytes() []byte {
	return b.bitmap
}

// Return number of bits in the bitmap
func (b *Bitmap) Len() int {
	return len(b.bitmap) * 8
}

func (b *Bitmap) String() string {
	var out string

	for _, byte_ := range b.bitmap {
		out += fmt.Sprintf("%08b ", byte_)
	}

	return out
}
