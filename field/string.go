package field

import (
	"encoding/json"
	"fmt"
)

var _ Field = (*String)(nil)
var _ json.Marshaler = (*String)(nil)
var _ json.Unmarshaler = (*String)(nil)

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
	if f.data != nil {
		*(f.data) = *f
	}
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
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	if len(packed) == 0 {
		return []byte{}, nil
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(packedLength, packed...), nil
}

func (f *String) Unpack(data []byte) (int, error) {
	dataLen, prefBytes, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}

	raw, read, err := f.spec.Enc.Decode(data[prefBytes:], dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %w", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	if err := f.SetBytes(raw); err != nil {
		return 0, fmt.Errorf("failed to set bytes: %w", err)
	}

	return read + prefBytes, nil
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

func (f *String) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal bytes to string: %w", err)
	}
	return f.SetBytes([]byte(v))
}
