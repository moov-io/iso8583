package iso8583

import (
	"errors"
	"reflect"
)

type Parser struct {
	messages  map[string]reflect.Type
	MtiEncode int
}

func (p *Parser) Register(mti string, tpl interface{}) {
	if len(mti) != 4 {
		panic("MTI must be a 4 digit numeric field")
	}
	v := reflect.ValueOf(tpl)
	// TODO do more check
	if p.messages == nil {
		p.messages = make(map[string]reflect.Type)
	}
	p.messages[mti] = reflect.Indirect(v).Type()
}

func decodeMti(raw []byte, encode int) (string, error) {
	mtiLen := 4
	if encode == BCD {
		mtiLen = 2
	}
	if len(raw) < mtiLen {
		return "", errors.New("bad raw data")
	}

	var mti string
	switch encode {
	case ASCII:
		mti = string(raw[:mtiLen])
	case BCD:
		mti = string(bcd2Ascii(raw[:mtiLen]))
	default:
		panic("invalid encode type")
	}
	return mti, nil
}

func (p *Parser) Parse(raw []byte) (*Message, error) {
	mti, err := decodeMti(raw, p.MtiEncode)
	if err != nil {
		return nil, err
	}

	tp, ok := p.messages[mti]
	if !ok {
		return nil, errors.New("no template registered for MTI: " + mti)
	}
	tpl := reflect.New(tp)
	initStruct(tp, tpl)
	msg := NewMessage(mti, tpl.Interface())
	msg.MtiEncode = p.MtiEncode
	return msg, msg.Load(raw)
}

func initStruct(tp reflect.Type, val reflect.Value) {
	for i := 0; i < tp.NumField(); i++ {
		field := reflect.Indirect(val).Field(i)
		fieldType := tp.Field(i)
		switch fieldType.Type.Kind() {
		case reflect.Ptr: // only initialize Ptr fields
			fieldValue := reflect.New(fieldType.Type.Elem())
			field.Set(fieldValue)
		}
	}
}
