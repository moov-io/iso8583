package encoding

var None Encoder = &noneEncoder{}

type noneEncoder struct{}

func (e *noneEncoder) Encode(data []byte) ([]byte, error) {
	return data, nil
}

func (e *noneEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	return data, 0, nil
}
