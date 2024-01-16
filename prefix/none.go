package prefix

var None = Prefixers{
	Fixed: &nonePrefixer{},
}

type nonePrefixer struct {
}

func (p *nonePrefixer) EncodeLength(int, int) ([]byte, error) {
	return []byte{}, nil
}

func (p *nonePrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return len(data), 0, nil
}

func (p *nonePrefixer) Inspect() string {
	return "None.Fixed"
}
