package field

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/moov-io/iso8583/utils"
)

var (
	_ Field            = (*String)(nil)
	_ json.Marshaler   = (*String)(nil)
	_ json.Unmarshaler = (*String)(nil)
)

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

	packer := f.spec.getPacker()

	return packer.Pack(data, f.spec)
}

func (f *String) Unpack(data []byte) (int, error) {
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
		case reflect.Int, reflect.Int64:
			i, err := strconv.Atoi(f.value)
			if err != nil {
				return fmt.Errorf("failed to convert string to int: %w", err)
			}

			val.SetInt(int64(i))
		case reflect.Slice:
			if val.Type().Elem().Kind() != reflect.Uint8 {
				return fmt.Errorf("can only be unmarshaled into []byte, got %s", val.Type())
			}
			val.SetBytes([]byte(f.value))
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
	case *int64:
		i, err := strconv.ParseInt(f.value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert string to int64: %w", err)
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
	if v == nil {
		f.value = ""
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.IsZero() {
		if !strings.Contains(reflect.ValueOf(v).Type().String(), "int") {
			f.value = ""
			return nil
		}
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
	case int64:
		f.value = strconv.FormatInt(v, 10)
	case *int:
		if v == nil {
			f.value = strconv.FormatInt(0, 10)
		} else {
			f.value = strconv.FormatInt(int64(*v), 10)
		}
	case *int64:
		if v == nil {
			f.value = strconv.FormatInt(0, 10)
		} else {
			f.value = strconv.FormatInt(*v, 10)
		}
	default:
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}

		//nolint:exhaustive
		switch kind {
		case reflect.String:
			f.value = rv.String()
		default:
			return fmt.Errorf("data does not match required *String or (string, *string, int, *int) type")
		}
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
