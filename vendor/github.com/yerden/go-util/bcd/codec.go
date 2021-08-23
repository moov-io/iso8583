package bcd

// Codec encapsulates both Encoder and Decoder.
type Codec struct {
	Encoder
	Decoder
}

// NewCodec returns new copy of Codec. See NewEncoder and NewDecoder
// on behaviour specifics.
func NewCodec(config *BCD) *Codec {
	return &Codec{*NewEncoder(config), *NewDecoder(config)}
}
