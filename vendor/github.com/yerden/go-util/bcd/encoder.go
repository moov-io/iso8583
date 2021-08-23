package bcd

// Encoder is used to encode decimal string into BCD bytes.
//
// Encoder may be copied with no side effects.
type Encoder struct {
	// symbol to nibble mapping; example:
	// '*' -> 0xA
	// the value > 0xf means no mapping, i.e. invalid symbol
	hash [0x100]byte

	// nibble used to fill if the number of bytes is odd
	filler byte

	// if true the 0x45 translates to '54' and vice versa
	swap bool
}

func checkBCD(config *BCD) bool {
	nibbles := make(map[byte]bool)
	// check all nibbles
	for _, nib := range config.Map {
		if _, ok := nibbles[nib]; ok || nib > 0xf {
			// already in map or not a nibble
			return false
		}
		nibbles[nib] = true
	}
	return config.Filler <= 0xf
}

func newHashEnc(config *BCD) (res [0x100]byte) {
	for i := 0; i < 0x100; i++ {
		c, ok := config.Map[byte(i)]
		if !ok {
			// no matching symbol
			c = 0xff
		}
		res[i] = c
	}
	return
}

// NewEncoder creates new Encoder from BCD configuration.  If the
// configuration is invalid NewEncoder will panic.
func NewEncoder(config *BCD) *Encoder {
	if !checkBCD(config) {
		panic("BCD table is incorrect")
	}
	return &Encoder{
		hash:   newHashEnc(config),
		filler: config.Filler,
		swap:   config.SwapNibbles}
}

func (enc *Encoder) packNibs(nib1, nib2 byte) byte {
	if enc.swap {
		return (nib2 << 4) + nib1&0xf
	} else {
		return (nib1 << 4) + nib2&0xf
	}
}

func (enc *Encoder) pack(w []byte) (n int, b byte, err error) {
	var nib1, nib2 byte
	switch len(w) {
	case 0:
		n = 0
		return
	case 1:
		n = 1
		if nib1, nib2 = enc.hash[w[0]], enc.filler; nib1 > 0xf {
			err = ErrBadInput
		}
	default:
		n = 2
		if nib1, nib2 = enc.hash[w[0]], enc.hash[w[1]]; nib1 > 0xf || nib2 > 0xf {
			err = ErrBadInput
		}
	}
	return n, enc.packNibs(nib1, nib2), err
}

// EncodedLen returns amount of space needed to store bytes after
// encoding data of length x.
func EncodedLen(x int) int {
	return (x + 1) / 2
}

// Encode get input bytes from src and encodes them into BCD data.
// Number of encoded bytes and possible error is returned.
func (enc *Encoder) Encode(dst, src []byte) (n int, err error) {
	var b byte
	var wid int

	for n < len(dst) {
		wid, b, err = enc.pack(src)
		switch {
		case err != nil:
			return
		case wid == 0:
			return
		}
		dst[n] = b
		n++
		src = src[wid:]
	}
	return
}
