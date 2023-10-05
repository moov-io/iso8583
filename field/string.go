package field

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*String)(nil)
var _ json.Marshaler = (*String)(nil)
var _ json.Unmarshaler = (*String)(nil)

type String struct {
	value string
	spec  *Spec
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

func (f *String) SetValue(v string) {
	f.value = v
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

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(data))
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

// Deprecated. Use Marshal instead
func (f *String) SetData(data interface{}) error {
	return f.Marshal(data)
}

func (f *String) Unmarshal(v interface{}) error {
	switch val := v.(type) {
	case reflect.Value:
		if !val.CanSet() {
			return fmt.Errorf("cannot set reflect.Value of type %s", val.Kind())
		}

		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			val.SetString(f.value)
		case reflect.Int:
			i, err := strconv.Atoi(f.value)
			if err != nil {
				return fmt.Errorf("failed to convert string to int: %w", err)
			}

			val.SetInt(int64(i))
		default:
			return fmt.Errorf("unsupported reflect.Value type: %s", val.Kind())
		}
	case *string:
		*val = f.value
	case *int:
		i, err := strconv.Atoi(f.value)
		if err != nil {
			return fmt.Errorf("failed to convert string to int: %w", err)
		}
		*val = i
	case *String:
		val.value = f.value
	default:
		return fmt.Errorf("unsupported type: expected *String, *string, or reflect.Value, got %T", v)
	}

	return nil
}

func (f *String) Marshal(v interface{}) error {
	if v == nil || (!reflect.ValueOf(v).CanInt() && reflect.ValueOf(v).IsZero()) {
		f.value = ""
		return nil
	}

	switch v := v.(type) {
	case *String:
		f.value = v.value
	case string:
		f.value = v
	case *string:
		f.value = *v
	case int:
		f.value = strconv.FormatInt(int64(v), 10)
	case *int:
		f.value = strconv.FormatInt(int64(*v), 10)
	default:
		return fmt.Errorf("data does not match required *String or (string, *string, int, *int) type")
	}

	return nil
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
