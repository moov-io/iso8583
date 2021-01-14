package field

import (
	"fmt"
)

var _ Field = (*String)(nil)

type String struct {
	Value string
	spec  *Spec
}

func NewString(spec *Spec) Field {
	return &String{
		spec: spec,
	}
}

func NewStringValue(val string) *String {
	return &String{
		Value: val,
	}
}

func (f *String) Spec() *Spec {
	return f.spec
}

func (f *String) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *String) SetBytes(b []byte) {
	f.Value = string(b)
}

func (f *String) Bytes() []byte {
	return []byte(f.Value)
}

func (f *String) String() string {
	return f.Value
}

func (f *String) Pack(data []byte) ([]byte, error) {
	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", f.spec.Description, err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", f.spec.Description, err)
	}

	return append(packedLength, packed...), nil
}

func (f *String) Unpack(data []byte) ([]byte, int, error) {
	dataLen, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	start := f.spec.Pref.Length()
	end := f.spec.Pref.Length() + dataLen
	raw, err := f.spec.Enc.Decode(data[start:end], f.spec.Length)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	f.Value = string(raw)

	return raw, dataLen + f.spec.Pref.Length(), nil
}
