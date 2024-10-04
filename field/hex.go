package field

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/moov-io/iso8583/utils"
)

var (
	_ Field            = (*Hex)(nil)
	_ json.Marshaler   = (*Hex)(nil)
	_ json.Unmarshaler = (*Hex)(nil)
)

// Hex field allows working with hex strings but under the hood it's a binary
// field. It's convenient to use when you need to work with hex strings, but
// don't want to deal with converting them to bytes manually.
// If provided value is not a valid hex string, it will return an error during
// packing. For the Hex field, the Binary encoding shoud be used in the Spec.
type Hex struct {
	value string
	spec  *Spec
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

	packer := f.spec.getPacker()

	return packer.Pack(data, f.spec)
}

func (f *Hex) Unpack(data []byte) (int, error) {
	unpacker := f.spec.getUnpacker()

	raw, bytesRead, err := unpacker.Unpack(data, f.spec)
	if err != nil {
		return 0, err
	}

	if err := f.SetBytes(raw); err != nil {
		return 0, fmt.Errorf("failed to set bytes: %w", err)
	}

	return bytesRead, nil
}

// Deprecated. Use Marshal instead
func (f *Hex) SetData(data interface{}) error {
	return f.Marshal(data)
}

func (f *Hex) Unmarshal(v interface{}) error {
	switch val := v.(type) {
	case reflect.Value:
		if !val.CanSet() {
			return fmt.Errorf("cannot set reflect.Value of type %s", val.Kind())
		}

		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			str, _ := f.String()
			val.SetString(str)
		case reflect.Slice:
			buf, _ := f.Bytes()
			val.SetBytes(buf)
		default:
			return fmt.Errorf("unsupported reflect.Value type: %s", val.Kind())
		}
	case *string:
		*val, _ = f.String()
	case *[]byte:
		*val, _ = f.Bytes()
	case *Hex:
		val.value = f.value
	default:
		return fmt.Errorf("unsupported type: expected *Hex, *string, *[]byte, or reflect.Value, got %T", v)
	}

	return nil
}

func (f *Hex) Marshal(v interface{}) error {
	if v == nil || reflect.ValueOf(v).IsZero() {
		f.value = ""
		return nil
	}

	switch v := v.(type) {
	case *Hex:
		f.value = v.value
	case string:
		f.value = v
	case *string:
		f.value = *v
	case []byte:
		f.SetBytes(v)
	case *[]byte:
		f.SetBytes(*v)
	default:
		return fmt.Errorf("data does not match required *Hex or (string, *string, []byte, *[]byte) type")
	}

	return nil
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
