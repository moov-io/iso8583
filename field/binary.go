package field

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/utils"
)

var (
	_ Field            = (*Binary)(nil)
	_ json.Marshaler   = (*Binary)(nil)
	_ json.Unmarshaler = (*Binary)(nil)
)

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

func (f *Binary) NewInstance() Field {
	return NewBinary(f.spec)
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

	packer := f.spec.getPacker()

	return packer.Pack(data, f.spec)
}

func (f *Binary) Unpack(data []byte) (int, error) {
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
			val.SetString(hex.EncodeToString(f.value))
		case reflect.Slice:
			if val.Type().Elem().Kind() != reflect.Uint8 {
				return fmt.Errorf("binary data can only be unmarshaled into []byte, got %s", val.Type())
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
	if v == nil {
		f.value = nil
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.IsZero() {
		f.value = nil
		return nil
	}

	switch v := v.(type) {
	case *Binary:
		f.value = v.value
	case string:
		buf, err := hex.DecodeString(v)
		if err != nil {
			return fmt.Errorf("failed to convert string to byte: %w", err)
		}

		f.value = buf
	case *string:
		buf, err := hex.DecodeString(*v)
		if err != nil {
			return fmt.Errorf("failed to convert string to byte: %w", err)
		}

		f.value = buf
	case []byte:
		f.SetBytes(v)
	case *[]byte:
		f.SetBytes(*v)
	default:
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}

		//nolint:exhaustive
		switch kind {
		case reflect.String:
			buf, err := hex.DecodeString(rv.String())
			if err != nil {
				return fmt.Errorf("failed to convert string to byte: %w", err)
			}

			f.value = buf
		default:
			return fmt.Errorf("data does not match required *Binary or (string, *string, []byte, *[]byte) type")
		}
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
