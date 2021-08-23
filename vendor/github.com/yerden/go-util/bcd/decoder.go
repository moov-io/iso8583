package bcd

type word [2]byte
type dword [4]byte
type qword [8]byte

// Decoder is used to decode BCD converted bytes into decimal string.
//
// Decoder may be copied with no side effects.
type Decoder struct {
	// if the input contains filler nibble in the middle, default
	// behaviour is to treat this as an error. You can tell decoder to
	// resume decoding quietly in that case by setting this.
	IgnoreFiller bool

	// two nibbles (1 byte) to 2 symbols mapping; example: 0x45 ->
	// '45' or '54' depending on nibble swapping additional 2 bytes of
	// dword should be 0, otherwise given byte is unacceptable
	hashWord [0x100]dword

	// one finishing byte with filler nibble to 1 symbol mapping;
	// example: 0x4f -> '4' (filler=0xf, swap=false)
	// additional byte of word should 0, otherise given nibble is
	// unacceptable
	hashByte [0x100]word
}

func newHashDecWord(config *BCD) (res [0x100]dword) {
	var w dword
	var b byte
	for i, _ := range res {
		// invalidating all bytes by default
		res[i] = dword{0xff, 0xff, 0xff, 0xff}
	}

	for c1, nib1 := range config.Map {
		for c2, nib2 := range config.Map {
			b = (nib1 << 4) + nib2&0xf
			if config.SwapNibbles {
				w = dword{c2, c1, 0, 0}
			} else {
				w = dword{c1, c2, 0, 0}
			}
			res[b] = w
		}
	}
	return
}

func newHashDecByte(config *BCD) (res [0x100]word) {
	var b byte
	for i, _ := range res {
		// invalidating all nibbles by default
		res[i] = word{0xff, 0xff}
	}
	for c, nib := range config.Map {
		if config.SwapNibbles {
			b = (config.Filler << 4) + nib&0xf
		} else {
			b = (nib << 4) + config.Filler&0xf
		}
		res[b] = word{c, 0}
	}
	return
}

func (dec *Decoder) unpack(w []byte, b byte) (n int, end bool, err error) {
	if dw := dec.hashWord[b]; dw[2] == 0 {
		return copy(w, dw[:2]), false, nil
	}
	if dw := dec.hashByte[b]; dw[1] == 0 {
		return copy(w, dw[:1]), true, nil
	}
	return 0, false, ErrBadBCD
}

// NewDecoder creates new Decoder from BCD configuration. If the
// configuration is invalid NewDecoder will panic.
func NewDecoder(config *BCD) *Decoder {
	if !checkBCD(config) {
		panic("BCD table is incorrect")
	}

	return &Decoder{
		hashWord: newHashDecWord(config),
		hashByte: newHashDecByte(config)}
}

// DecodedLen tells how much space is needed to store decoded string.
// Please note that it returns the max amount of possibly needed space
// because last octet may contain only one encoded digit. In that
// case the decoded length will be less by 1. For example, 4 octets
// may encode 7 or 8 digits.  Please examine the result of Decode to
// obtain the real value.
func DecodedLen(x int) int {
	return 2 * x
}

// Decode parses BCD encoded bytes from src and tries to decode them
// to dst. Number of decoded bytes and possible error is returned.
func (dec *Decoder) Decode(dst, src []byte) (n int, err error) {
	if len(src) == 0 {
		return 0, nil
	}

	for _, c := range src[:len(src)-1] {
		wid, end, err := dec.unpack(dst[n:], c)
		switch {
		case err != nil: // invalid input
			return n, err
		case wid == 0: // no place in dst
			return n, nil
		case end && !dec.IgnoreFiller: // unexpected filler
			return n, ErrBadBCD
		}
		n += wid
	}

	c := src[len(src)-1]
	wid, _, err := dec.unpack(dst[n:], c)
	switch {
	case err != nil: // invalid input
		return n, err
	case wid == 0: // no place in dst
		return n, nil
	}
	n += wid
	return n, nil
}
