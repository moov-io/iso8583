package field

import (
	"encoding/json"
	"fmt"
)

var _ Field = (*String)(nil)

type String struct {
	Value string `json:"value"`
	spec  *Spec
	data  *String
}

func NewString(spec *Spec) *String {
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

func (f *String) SetBytes(b []byte) error {
	f.Value = string(b)
	return nil
}

func (f *String) Bytes() ([]byte, error) {
	return []byte(f.Value), nil
}

func (f *String) String() (string, error) {
	return f.Value, nil
}

func (f *String) Pack() ([]byte, error) {
	data := []byte(f.Value)

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %v", err)
	}

	valueLen := len(packed)
	if f.spec.CountT == "1" {
		valueLen = len(data)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, valueLen)
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %v", err)
	}

	return append(packedLength, packed...), nil
}

func (f *String) Unpack(data []byte) (int, error) {
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

	f.Value = string(raw)

	if f.data != nil {
		*(f.data) = *f
	}

	return read + f.spec.Pref.Length(), nil
}

func (f *String) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	str, ok := data.(*String)
	if !ok {
		return fmt.Errorf("data does not match required *String type")
	}

	f.data = str
	if str.Value != "" {
		f.Value = str.Value
	}
	return nil
}

func (f *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}
