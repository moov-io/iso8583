package field

import (
	"fmt"
	"strconv"
)

var _ Field = (*Numeric)(nil)

type Numeric struct {
	Value int
	spec  *Spec
}

func NewNumeric(spec *Spec) Field {
	return &Numeric{
		spec: spec,
	}
}

func NewNumericValue(val int) *Numeric {
	return &Numeric{
		Value: val,
	}
}

func (f *Numeric) Spec() *Spec {
	return f.spec
}

func (f *Numeric) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Numeric) SetBytes(b []byte) {
	f.Value, _ = strconv.Atoi(string(b))
}

func (f *Numeric) Bytes() []byte {
	return []byte(strconv.Itoa(f.Value))
}

func (f *Numeric) String() string {
	return strconv.Itoa(f.Value)
}

func (f *Numeric) Pack() ([]byte, error) {
	data := []byte(strconv.Itoa(f.Value))

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %v", err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %v", err)
	}

	return append(packedLength, packed...), nil
}

// returns number of bytes was read
func (f *Numeric) Unpack(data []byte) (int, error) {
	dataLen, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %v", err)
	}

	start := f.spec.Pref.Length()
	raw, read, err := f.spec.Enc.Decode(data[start:], dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %v", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	f.Value, err = strconv.Atoi(string(raw))
	if err != nil {
		return 0, fmt.Errorf("failed to convert into number: %v", err)
	}

	return read + f.spec.Pref.Length(), nil
}
