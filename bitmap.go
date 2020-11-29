package iso8583

// Partially copied from https://github.com/boljen/go-bitmap (MIT license)

import (
	"fmt"

	"github.com/boljen/go-bitmap"
)

type Bitmap struct {
	bitmap bitmap.Bitmap
}

func NewBitmap(len int) *Bitmap {
	return &Bitmap{
		bitmap: bitmap.New(len),
	}
}

func (b *Bitmap) IsSet(i int) bool {
	return b.bitmap.Get(i)
}

func (b *Bitmap) Set(i int) {
	b.bitmap.Set(i, true)
}

func (b *Bitmap) Bytes() []byte {
	return []byte(b.bitmap)
}

func (b *Bitmap) String() string {
	var out string

	for _, byte_ := range b.bitmap {
		out += fmt.Sprintf("%08b ", byte_)
	}

	return out
}
