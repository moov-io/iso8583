package encoding

var Binary Encoder = &binaryEncoder{}

type binaryEncoder struct{}

func (e binaryEncoder) Encode(data []byte) ([]byte, error) {
	out := append([]byte(nil), data...)

	return out, nil
}

func (e binaryEncoder) Decode(data []byte) ([]byte, error) {
	out := append([]byte(nil), data...)

	return out, nil
}
