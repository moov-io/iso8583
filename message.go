package iso8583

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	TAG_FIELD  string = "field"
	TAG_ENCODE string = "encode"
	TAG_LENGTH string = "length"
)

type fieldInfo struct {
	Index     int
	Encode    int
	LenEncode int
	Length    int
	Field     Iso8583Type
}

// Message is structure for ISO 8583 message encode and decode
type Message struct {
	Mti          string
	MtiEncode    int
	SecondBitmap bool
	Data         interface{}
}

// NewMessage creates new Message structure
func NewMessage(mti string, data interface{}) *Message {
	return &Message{mti, ASCII, false, data}
}

// Bytes marshall Message to bytes
func (m *Message) Bytes() (ret []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Critical error:" + fmt.Sprint(r))
			ret = nil
		}
	}()

	ret = make([]byte, 0)

	// generate MTI:
	mtiBytes, err := m.encodeMti()
	if err != nil {
		return nil, err
	}
	ret = append(ret, mtiBytes...)

	// generate bitmap and fields:
	fields := parseFields(m.Data)

	byteNum := 8
	if m.SecondBitmap {
		byteNum = 16
	}
	bitmap := make([]byte, byteNum)
	data := make([]byte, 0, 512)

	for byteIndex := 0; byteIndex < byteNum; byteIndex++ {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {

			i := byteIndex*8 + bitIndex + 1

			// if we need second bitmap (additional 8 bytes) - set first bit in first bitmap
			if m.SecondBitmap && i == 1 {
				step := uint(7 - bitIndex)
				bitmap[byteIndex] |= (0x01 << step)
			}

			if info, ok := fields[i]; ok {

				// if field is empty, then we can't add it to bitmap
				if info.Field.IsEmpty() {
					continue
				}

				// mark 1 in bitmap:
				step := uint(7 - bitIndex)
				bitmap[byteIndex] |= (0x01 << step)
				// append data:
				d, err := info.Field.Bytes(info.Encode, info.LenEncode, info.Length)
				if err != nil {
					return nil, err
				}
				data = append(data, d...)
			}
		}
	}
	ret = append(ret, bitmap...)
	ret = append(ret, data...)

	return ret, nil
}

func (m *Message) encodeMti() ([]byte, error) {
	if m.Mti == "" {
		return nil, errors.New("MTI is required")
	}
	if len(m.Mti) != 4 {
		return nil, errors.New("MTI is invalid")
	}

	switch m.MtiEncode {
	case BCD:
		return bcd([]byte(m.Mti)), nil
	default:
		return []byte(m.Mti), nil
	}
}

func parseFields(msg interface{}) map[int]*fieldInfo {
	fields := make(map[int]*fieldInfo)

	v := reflect.Indirect(reflect.ValueOf(msg))
	if v.Kind() != reflect.Struct {
		panic("data must be a struct")
	}
	for i := 0; i < v.NumField(); i++ {
		if isPtrOrInterface(v.Field(i).Kind()) && v.Field(i).IsNil() {
			continue
		}

		sf := v.Type().Field(i)
		if sf.Tag == "" || sf.Tag.Get(TAG_FIELD) == "" {
			continue
		}

		index, err := strconv.Atoi(sf.Tag.Get(TAG_FIELD))
		if err != nil {
			panic("value of field must be numeric")
		}

		encode := 0
		lenEncode := 0
		if raw := sf.Tag.Get(TAG_ENCODE); raw != "" {
			enc := strings.Split(raw, ",")
			if len(enc) == 2 {
				lenEncode = parseEncodeStr(enc[0])
				encode = parseEncodeStr(enc[1])
			} else {
				encode = parseEncodeStr(enc[0])
			}
		}

		length := -1
		if l := sf.Tag.Get(TAG_LENGTH); l != "" {
			length, err = strconv.Atoi(l)
			if err != nil {
				panic("value of length must be numeric")
			}
		}

		field, ok := v.Field(i).Interface().(Iso8583Type)
		if !ok {
			panic("field must be Iso8583Type")
		}
		fields[index] = &fieldInfo{index, encode, lenEncode, length, field}
	}
	return fields
}

func isPtrOrInterface(k reflect.Kind) bool {
	return k == reflect.Interface || k == reflect.Ptr
}

func parseEncodeStr(str string) int {
	switch str {
	case "ascii":
		return ASCII
	case "lbcd":
		fallthrough
	case "bcd":
		return BCD
	case "rbcd":
		return rBCD
	}
	return -1
}

// Load unmarshall Message from bytes
func (m *Message) Load(raw []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Critical error:" + fmt.Sprint(r))
		}
	}()

	if m.Mti == "" {
		m.Mti, err = decodeMti(raw, m.MtiEncode)
		if err != nil {
			return err
		}
	}
	start := 4
	if m.MtiEncode == BCD {
		start = 2
	}

	fields := parseFields(m.Data)

	byteNum := 8
	if raw[start]&0x80 == 0x80 {
		// 1st bit == 1
		m.SecondBitmap = true
		byteNum = 16
	}
	bitByte := raw[start : start+byteNum]
	start += byteNum

	for byteIndex := 0; byteIndex < byteNum; byteIndex++ {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			step := uint(7 - bitIndex)
			if (bitByte[byteIndex] & (0x01 << step)) == 0 {
				continue
			}

			i := byteIndex*8 + bitIndex + 1
			if i == 1 {
				// field 1 is the second bitmap
				continue
			}
			f, ok := fields[i]
			if !ok {
				return fmt.Errorf("field %d not defined", i)
			}
			l, err := f.Field.Load(raw[start:], f.Encode, f.LenEncode, f.Length)
			if err != nil {
				return err
			}
			start += l
		}
	}
	return nil
}
