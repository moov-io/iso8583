package padding

type Padder interface {
	Pad(data []byte, length int) []byte
	Unpad(data []byte) []byte
	Inspect() []byte
}
