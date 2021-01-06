package field

import (
	"fmt"
)

var _ Field = (*StringField)(nil)

type StringField struct {
	Value string
	spec  *Spec
}

func NewStringField(spec *Spec) Field {
	return &StringField{
		spec: spec,
	}
}

func NewStringValue(val string) *StringField {
	return &StringField{
		Value: val,
	}
}

func (f *StringField) Spec() *Spec {
	return f.spec
}

func (f *StringField) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *StringField) SetBytes(b []byte) {
	f.Value = string(b)
}

func (f *StringField) Bytes() []byte {
	return []byte(f.Value)
}

func (f *StringField) String() string {
	return f.Value
}

func (f *StringField) Pack(data []byte) ([]byte, error) {
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

func (f *StringField) Unpack(data []byte) ([]byte, int, error) {
	dataLen, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	start := f.spec.Pref.Length()
	end := f.spec.Pref.Length() + dataLen
	raw, err := f.spec.Enc.Decode(data[start:end])
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	f.Value = string(raw)

	return raw, dataLen + f.spec.Pref.Length(), nil
}
