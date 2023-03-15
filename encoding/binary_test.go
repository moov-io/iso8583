package encoding

import "testing"

func FuzzDecodeBinary(f *testing.F) {
	enc := &binaryEncoder{}

	f.Fuzz(func(t *testing.T, data []byte, length int) {
		enc.Decode(data, length)
	})
}
