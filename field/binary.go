package field

import (
	"encoding/json"
	"fmt"

	"github.com/moov-io/iso8583/encoding"
)

var _ Field = (*Binary)(nil)
var _ json.Marshaler = (*Binary)(nil)
var _ json.Unmarshaler = (*Binary)(nil)

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
	if f.data != nil {
		*(f.data) = *f
	}
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

func (f *Binary) Unpack(data []byte) (int, error) {
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
	str, err := f.String()
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

func (f *Binary) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal bytes to string: %w", err)
	}

	hex, err := encoding.ASCIIHexToBytes.Encode([]byte(v))
	if err != nil {
		return fmt.Errorf("failed to convert ASCII Hex string to bytes")
	}
	return f.SetBytes(hex)
}
