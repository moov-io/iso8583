package field

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*Hex)(nil)
var _ json.Marshaler = (*Hex)(nil)
var _ json.Unmarshaler = (*Hex)(nil)

// Hex is a field that contains a hex string value, but is encoded as binary
type Hex struct {
	value string
	spec  *Spec
	data  *Hex
}

func NewHex(spec *Spec) *Hex {
	return &Hex{
		spec: spec,
	}
}

// NewHexValue creates a new Hex field with the given value The value is
// converted from hex to bytes before packing, so we don't validate that val is
// a valid hex string here.
func NewHexValue(val string) *Hex {
	return &Hex{
		value: val,
	}
}

func (f *Hex) Spec() *Spec {
	return f.spec
}

func (f *Hex) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Hex) SetBytes(b []byte) error {
	f.value = strings.ToUpper(hex.EncodeToString(b))
	if f.data != nil {
		*(f.data) = *f
	}
	return nil
}

func (f *Hex) Bytes() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return hex.DecodeString(f.value)
}

func (f *Hex) String() (string, error) {
	if f == nil {
		return "", nil
	}
	return f.value, nil
}

func (f *Hex) Value() string {
	if f == nil {
		return ""
	}
	return f.value
}

func (f *Hex) SetValue(v string) {
	f.value = v
}

func (f *Hex) Pack() ([]byte, error) {
	data, err := f.Bytes()
	if err != nil {
		return nil, utils.NewSafeErrorf(err, "converting hex field into bytes")
	}

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(packedLength, packed...), nil
}

func (f *Hex) Unpack(data []byte) (int, error) {
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

func (f *Hex) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	str, ok := v.(*Hex)
	if !ok {
		return errors.New("data does not match required *Hex type")
	}

	str.value = f.value

	return nil
}

func (f *Hex) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	str, ok := data.(*Hex)
	if !ok {
		return fmt.Errorf("data does not match required *Hex type")
	}

	f.data = str
	if str.value != "" {
		f.value = str.value
	}
	return nil
}

func (f *Hex) Marshal(data interface{}) error {
	return f.SetData(data)
}

func (f *Hex) MarshalJSON() ([]byte, error) {
	data, err := f.String()
	if err != nil {
		return nil, utils.NewSafeError(err, "convert hex field into bytes")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal string to bytes")
	}
	return bytes, nil
}

func (f *Hex) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return utils.NewSafeError(err, "failed to JSON unmarshal bytes to string")
	}

	f.value = v

	return nil
}
