package field

import (
	"encoding/json"
	"fmt"
)

var _ Field = (*Binary)(nil)

type Binary struct {
	Value []byte `json:"value"`
	spec  *Spec
	data  *Binary
}

func NewBinary(spec *Spec) *Binary {
	return &Binary{
		spec: spec,
	}
}

func NewBinaryValue(val []byte) *Binary {
	return &Binary{
		Value: val,
	}
}

func (f *Binary) Spec() *Spec {
	return f.spec
}

func (f *Binary) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Binary) SetBytes(b []byte) error {
	f.Value = b
	return nil
}

func (f *Binary) Bytes() ([]byte, error) {
	return f.Value, nil
}

func (f *Binary) String() (string, error) {
	return fmt.Sprintf("%X", f.Value), nil
}

func (f *Binary) Pack() ([]byte, error) {
	data := f.Value

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

func (f *Binary) Unpack(data []byte) (int, error) {
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

	f.Value = raw

	if f.data != nil {
		*(f.data) = *f
	}

	return read + f.spec.Pref.Length(), nil
}

func (f *Binary) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	str, ok := data.(*Binary)
	if !ok {
		return fmt.Errorf("data does not match required *Binary type")
	}

	f.data = str
	if str.Value != nil {
		f.Value = str.Value
	}
	return nil
}

func (f *Binary) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}
