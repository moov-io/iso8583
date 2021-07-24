package padding

var None Padder = &nonePadder{}

type nonePadder struct{}

func (p *nonePadder) Pad(data []byte, length int) []byte {
	return data
}

func (p *nonePadder) Unpad(data []byte) []byte {
	return data
}

func (p *nonePadder) Inspect() []byte {
	return nil
}
