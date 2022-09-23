package field

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*String)(nil)
var _ json.Marshaler = (*String)(nil)
var _ json.Unmarshaler = (*String)(nil)

type String struct {
	value string
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
		value: val,
	}
}

func (f *String) Spec() *Spec {
	return f.spec
}

func (f *String) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *String) SetBytes(b []byte) error {
	f.value = string(b)
	if f.data != nil {
		*(f.data) = *f
	}
	return nil
}

func (f *String) Bytes() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return []byte(f.value), nil
}

func (f *String) String() (string, error) {
	if f == nil {
		return "", nil
	}
	return f.value, nil
}

func (f *String) Value() string {
	if f == nil {
		return ""
	}
	return f.value
}

func (f *String) Pack() ([]byte, error) {
	data := []byte(f.value)

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
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

func (f *String) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	str, ok := v.(*String)
	if !ok {
		return errors.New("data does not match required *String type")
	}

	str.value = f.value

	return nil
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
	if str.value != "" {
		f.value = str.value
	}
	return nil
}

func (f *String) Marshal(data interface{}) error {
	return f.SetData(data)
}

func (f *String) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(f.value)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal string to bytes")
	}
	return bytes, nil
}

func (f *String) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return utils.NewSafeError(err, "failed to JSON unmarshal bytes to string")
	}
	return f.SetBytes([]byte(v))
}
