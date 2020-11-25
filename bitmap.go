package iso8583

import "fmt"

// source https://play.golang.org/p/oOe1Gd4C2G
type Bitmap [8]byte

func (bits *Bitmap) IsSet(i int) bool {
	i -= 1
	return bits[i/8]&(1<<uint(7-i%8)) != 0
}

func (bits *Bitmap) Set(i int) {
	i -= 1
	bits[i/8] |= 1 << uint(7-i%8)
}

func (bits *Bitmap) Clear(i int) {
	i -= 1
	bits[i/8] &^= 1 << uint(7-i%8)
}

func (bits *Bitmap) String() string {
	var out string

	for _, b := range bits {
		out += fmt.Sprintf("%08b ", b)
	}

	return out
}
