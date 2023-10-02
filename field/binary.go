package field

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*Binary)(nil)
var _ json.Marshaler = (*Binary)(nil)
var _ json.Unmarshaler = (*Binary)(nil)

type Binary struct {
	value []byte
	spec  *Spec
}

func NewBinary(spec *Spec) *Binary {
	return &Binary{
		spec: spec,
	}
}

func NewBinaryValue(val []byte) *Binary {
	return &Binary{
		value: val,
	}
}

func (f *Binary) Spec() *Spec {
	return f.spec
}

func (f *Binary) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Binary) SetBytes(b []byte) error {
	f.value = b
	return nil
}

func (f *Binary) Bytes() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return f.value, nil
}

func (f *Binary) String() (string, error) {
	if f == nil {
		return "", nil
	}
	return fmt.Sprintf("%X", f.value), nil
}

func (f *Binary) Value() []byte {
	if f == nil {
		return nil
	}
	return f.value
}

func (f *Binary) SetValue(v []byte) {
	f.value = v
}

func (f *Binary) Pack() ([]byte, error) {
	data := f.value

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

// Deprecated. Use Marshal instead
func (f *Binary) SetData(data interface{}) error {
	return f.Marshal(data)
}

func (f *Binary) Unmarshal(v interface{}) error {
	switch val := v.(type) {
	case reflect.Value:
		if !val.CanSet() {
			return fmt.Errorf("cannot set reflect.Value of type %s", val.Kind())
		}

		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			if !val.CanSet() {
				return fmt.Errorf("reflect.Value of the data can not be change")
			}

			str := hex.EncodeToString(f.value)
			val.SetString(str)
		case reflect.Slice:
			if !val.CanSet() {
				return fmt.Errorf("reflect.Value of the data can not be change")
			}

			val.SetBytes(f.value)
		default:
			return fmt.Errorf("unsupported reflect.Value type: %s", val.Kind())
		}
	case *string:
		*val = hex.EncodeToString(f.value)
	case *[]byte:
		*val = f.value
	case *Binary:
		val.value = f.value
	default:
		return fmt.Errorf("unsupported type: expected *Binary, *string, *[]byte, or reflect.Value, got %T", v)
	}

	return nil
}

func (f *Binary) Marshal(v interface{}) error {
	switch v := v.(type) {
	case *Binary:
		if v == nil {
			return nil
		}
		f.value = v.value
	case string:
		if v == "" {
			f.value = nil
			return nil
		}

		buf, err := hex.DecodeString(v)
		if err != nil {
			return fmt.Errorf("failed to convert string to byte: %w", err)
		}

		f.value = buf
	case *string:
		if v == nil {
			f.value = nil
			return nil
		}

		buf, err := hex.DecodeString(*v)
		if err != nil {
			return fmt.Errorf("failed to convert string to byte: %w", err)
		}

		f.value = buf
	case []byte:
		f.SetBytes(v)
	case *[]byte:
		if v == nil {
			f.value = nil
			return nil
		}
		f.SetBytes(*v)
	default:
		return fmt.Errorf("data does not match required *Binary or (string, *string, []byte, *[]byte) type")
	}

	return nil
}

func (f *Binary) MarshalJSON() ([]byte, error) {
	str, err := f.String()
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(str)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal string to bytes")
	}
	return bytes, nil
}

func (f *Binary) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return utils.NewSafeError(err, "failed to JSON unmarshal bytes to string")
	}

	hex, err := encoding.ASCIIHexToBytes.Encode([]byte(v))
	if err != nil {
		return utils.NewSafeError(err, "failed to convert ASCII Hex string to bytes")
	}
	return f.SetBytes(hex)
}
