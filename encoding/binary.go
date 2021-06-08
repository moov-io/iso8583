package encoding

var Binary Encoder = &binaryEncoder{}

type binaryEncoder struct{}

func (e binaryEncoder) Encode(data []byte) ([]byte, error) {
	out := append([]byte(nil), data...)

	return out, nil
}

func (e binaryEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	out := append([]byte(nil), data...)

	return out[:length], length, nil
}
