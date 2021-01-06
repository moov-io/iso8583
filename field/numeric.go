package field

import (
	"fmt"
	"strconv"
)

var _ Field = (*NumericField)(nil)

type NumericField struct {
	Value int
	spec  *Spec
}

func NewNumericField(spec *Spec) Field {
	return &NumericField{
		spec: spec,
	}
}

func NewNumericValue(val int) *NumericField {
	return &NumericField{
		Value: val,
	}
}

func (f *NumericField) Spec() *Spec {
	return f.spec
}

func (f *NumericField) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *NumericField) SetBytes(b []byte) {
	f.Value, _ = strconv.Atoi(string(b))
}

func (f *NumericField) Bytes() []byte {
	return []byte(strconv.Itoa(f.Value))
}

func (f *NumericField) String() string {
	return strconv.Itoa(f.Value)
}

func (f *NumericField) Pack(data []byte) ([]byte, error) {
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

func (f *NumericField) Unpack(data []byte) ([]byte, int, error) {
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

	f.Value, _ = strconv.Atoi(string(raw))

	return raw, dataLen + f.spec.Pref.Length(), nil
}
