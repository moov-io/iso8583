package iso8583

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DataIso struct {
	F2  *Llnumeric `field:"2" length:"19"`
	F3  *Numeric   `field:"3" length:"6"`
	F4  *Numeric   `field:"4" length:"12"`
	F7  *Numeric   `field:"7" length:"10"`
	F11 *Numeric   `field:"11" length:"6"`
	F12 *Numeric   `field:"12" length:"6"`
	F13 *Numeric   `field:"13" length:"4"`
	F14 *Numeric   `field:"14" length:"4"`
	// BCD encoding with right-aligned value with odd length (for ex. "643" as [6 67] == "0643")
	F19  *Numeric      `field:"19" length:"3" encode:"rbcd"`
	F22  *Numeric      `field:"22" length:"3"`
	F25  *Numeric      `field:"25" length:"2"`
	F32  *Llnumeric    `field:"32" length:"11"`
	F35  *Llnumeric    `field:"35" length:"37"`
	F37  *Alphanumeric `field:"37" length:"12"`
	F39  *Alphanumeric `field:"39" length:"2"`
	F41  *Alphanumeric `field:"41" length:"8"`
	F42  *Alphanumeric `field:"42" length:"15"`
	F43  *Alphanumeric `field:"43" length:"40"`
	F49  *Numeric      `field:"49" length:"3" encode:"bcd"`
	F52  *Binary       `field:"52" length:"8"`
	F53  *Numeric      `field:"53" length:"16"`
	F120 *Lllnumeric   `field:"120" length:"999"`
}

func TestEncode(t *testing.T) {
	data := &DataIso{
		F2:   &Llnumeric{"4276555555555555"},
		F3:   &Numeric{"000000"},
		F4:   &Numeric{"000000077700"},
		F7:   &Numeric{"0701111844"},
		F11:  &Numeric{"000123"},
		F12:  &Numeric{"131844"},
		F13:  &Numeric{"0701"},
		F14:  &Numeric{"1902"},
		F19:  &Numeric{"643"},
		F22:  &Numeric{"901"},
		F25:  &Numeric{"02"},
		F32:  &Llnumeric{"123456"},
		F35:  &Llnumeric{"4276555555555555=12345678901234567890"},
		F37:  &Alphanumeric{"987654321001"},
		F39:  NewAlphanumeric(""),
		F41:  &Alphanumeric{"00000321"},
		F42:  &Alphanumeric{"120000000000034"},
		F43:  &Alphanumeric{"Test text"},
		F49:  &Numeric{"643"},
		F52:  NewBinary([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
		F53:  &Numeric{"1234000000000000"},
		F120: &Lllnumeric{"Another test text"},
	}

	iso := Message{"0100", ASCII, true, data}

	res, err := iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	sample := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	if bytes.Compare(res, sample) != 0 {
		t.Error("ISO Encode error!")
	}
}

func TestDecode(t *testing.T) {

	input := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	// init empty iso message struct
	iso := Message{"", ASCII, true, newDataIso()}

	// parse data from bytes to iso struct
	err := iso.Load(input)

	if err != nil {
		t.Error("ISO Decode error:", err)
	}

	resultFields := iso.Data.(*DataIso)

	// check BCD numeric values length
	assert.Equal(t, 3, len(resultFields.F19.Value))
	assert.Equal(t, 3, len(resultFields.F49.Value))

	// check values for BCD (lBCD) and rBCD
	assert.Equal(t, "643", resultFields.F19.Value)
	assert.Equal(t, "643", resultFields.F49.Value)

	var res []byte

	// set second bitmap because field 120 in struct (need if more than 63 fields in message)
	iso.SecondBitmap = true

	// before encode add "0" to left of F19 for testing rBCD encoding
	iso.Data.(*DataIso).F19.Value = "0" + iso.Data.(*DataIso).F19.Value

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	// set field 120 value to empty string
	iso.Data.(*DataIso).F120.Value = ""

	iso.SecondBitmap = false

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	sample := []byte{48, 49, 48, 48, 114, 60, 36, 129, 40, 224, 152, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}

	if bytes.Compare(res, sample) != 0 {
		t.Error("ISO Encode error!")
	}
}

func TestParserErrors(t *testing.T) {

	parser := Parser{}

	err := parser.Register("0100", nil)

	assert.EqualError(t, err, "Critical error:reflect: call of reflect.Value.Type on zero Value")

	err = parser.Register("1", newDataIso())

	assert.EqualError(t, err, "MTI must be a 4 digit numeric field")

	_, err = parser.Parse([]byte{0})

	assert.EqualError(t, err, "bad raw data")

	parser.MtiEncode = BCD

	_, err = parser.Parse([]byte{1, 2})

	assert.EqualError(t, err, "no template registered for MTI: 0102")

	parser.MtiEncode = 10

	_, err = parser.Parse([]byte{1, 2, 3, 4})

	assert.EqualError(t, err, "invalid encode type")

	parser.MtiEncode = ASCII

	input := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	err = parser.Register("0100", newDataIso())

	_, err = parser.Parse(input[0:23])

	assert.EqualError(t, err, "Critical error:runtime error: slice bounds out of range")

	parser.messages["0100"] = nil

	_, err = parser.Parse(input)

	assert.EqualError(t, err, "Critical error:reflect: New(nil)")
}

func TestParser(t *testing.T) {

	input := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	parser := Parser{}

	err := parser.Register("0100", newDataIso())

	assert.Equal(t, nil, err)
	// parse data from bytes to iso struct
	// parse data from bytes to iso struct
	iso, err := parser.Parse(input)

	if err != nil {
		t.Error("ISO Decode error:", err)
	}

	resultFields := iso.Data.(*DataIso)

	// check BCD numeric values length
	assert.Equal(t, 3, len(resultFields.F19.Value))
	assert.Equal(t, 3, len(resultFields.F49.Value))

	// check values for BCD (lBCD) and rBCD
	assert.Equal(t, "643", resultFields.F19.Value)
	assert.Equal(t, "643", resultFields.F49.Value)

	var res []byte

	// set second bitmap because field 120 in struct (need if more than 63 fields in message)
	iso.SecondBitmap = true

	// before encode add "0" to left of F19 for testing rBCD encoding
	iso.Data.(*DataIso).F19.Value = "0" + iso.Data.(*DataIso).F19.Value

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	// set field 120 value to empty string
	iso.Data.(*DataIso).F120.Value = ""

	iso.SecondBitmap = false

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	sample := []byte{48, 49, 48, 48, 114, 60, 36, 129, 40, 224, 152, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}

	if bytes.Compare(res, sample) != 0 {
		t.Error("ISO Encode error!")
	}
}

func TestMessage(t *testing.T) {
	type TestIso struct {
		DataIso
		AB *Llnumeric `field:"ab" length:"19"`
	}

	iso := Message{"", ASCII, true, TestIso{*newDataIso(), NewLlnumeric("")}}

	input := []byte{48, 49, 48, 48, 114, 60, 36, 129, 40, 224, 152, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}

	err := iso.Load(input)

	assert.EqualError(t, err, "Critical error:value of field must be numeric")

	type TestIso2 struct {
		F2 *Llnumeric `field:"2" length:"19"`
	}

	iso = Message{"", ASCII, true, TestIso2{}}

	err = iso.Load(input)

	assert.EqualError(t, err, "field 2 not defined")

}

func TestBCD(t *testing.T) {

	b := []byte("954")
	r := rbcd(b)
	assert.Equal(t, "0954", fmt.Sprintf("%X", r))

	r = lbcd(b)
	assert.Equal(t, "9540", fmt.Sprintf("%X", r))

	b = []byte("31")
	r = lbcd(b)
	assert.Equal(t, "31", fmt.Sprintf("%X", r))
	r = rbcd(b)
	assert.Equal(t, "31", fmt.Sprintf("%X", r))

	assert.Panics(t,
		func() {
			bcd([]byte{0})
		}, "Calling bcd() with len(data) % 2 != 0 should panic")

}

// newDataIso creates DataIso
func newDataIso() *DataIso {
	return &DataIso{
		F2:   NewLlnumeric(""),
		F3:   NewNumeric(""),
		F4:   NewNumeric(""),
		F7:   NewNumeric(""),
		F11:  NewNumeric(""),
		F12:  NewNumeric(""),
		F13:  NewNumeric(""),
		F14:  NewNumeric(""),
		F19:  NewNumeric(""),
		F22:  NewNumeric(""),
		F25:  NewNumeric(""),
		F32:  NewLlnumeric(""),
		F35:  NewLlnumeric(""),
		F37:  NewAlphanumeric(""),
		F39:  NewAlphanumeric(""),
		F41:  NewAlphanumeric(""),
		F42:  NewAlphanumeric(""),
		F43:  NewAlphanumeric(""),
		F49:  NewNumeric(""),
		F52:  NewBinary(nil),
		F53:  NewNumeric(""),
		F120: NewLllnumeric(""),
	}
}
