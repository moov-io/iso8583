package iso8583

import (
	"testing"
	"bytes"
)

type DataIso struct {
	F2   *Llnumeric     `field:"2" length:"19"`
	F3   *Numeric       `field:"3" length:"6"`
	F4   *Numeric       `field:"4" length:"12"`
	F7   *Numeric       `field:"7" length:"10"`
	F11  *Numeric       `field:"11" length:"6"`
	F12  *Numeric       `field:"12" length:"6"`
	F13  *Numeric       `field:"13" length:"4"`
	F14  *Numeric       `field:"14" length:"4"`
	F19  *Numeric       `field:"19" length:"3"`
	F22  *Numeric       `field:"22" length:"3"`
	F25  *Numeric       `field:"25" length:"2"`
	F32  *Llnumeric     `field:"32" length:"11"`
	F35  *Llnumeric     `field:"35" length:"37"`
	F37  *Alphanumeric  `field:"37" length:"12"`
	F39  *Alphanumeric  `field:"39" length:"2"`
	F41  *Alphanumeric  `field:"41" length:"8"`
	F42  *Alphanumeric  `field:"42" length:"15"`
	F43  *Alphanumeric  `field:"43" length:"40"`
	F49  *Alphanumeric  `field:"49" length:"3"`
	F52  *Binary        `field:"52" length:"8"`
	F53  *Numeric       `field:"53" length:"16"`
	F120 *Lllnumeric    `field:"120" length:"999"`
}

func TestEncode(t *testing.T) {
	data := &DataIso{
		F2 : &Llnumeric{"4276555555555555"},
		F3 : &Numeric{"000000"},
		F4 : &Numeric{"000000077700"},
		F7 : &Numeric{"0701111844"},
		F11: &Numeric{"000123"},
		F12: &Numeric{"131844"},
		F13: &Numeric{"0701"},
		F14: &Numeric{"1902"},
		F19: &Numeric{"643"},
		F22: &Numeric{"901"},
		F25: &Numeric{"02"},
		F32: &Llnumeric{"123456"},
		F35: &Llnumeric{"4276555555555555=12345678901234567890"},
		F37: &Alphanumeric{"987654321001"},
		F39: &Alphanumeric{},
		F41: &Alphanumeric{"00000321"},
		F42: &Alphanumeric{"120000000000034"},
		F43: &Alphanumeric{"Test text"},
		F49: &Alphanumeric{"643"},
		F52: NewBinary([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
		F53: &Numeric{"1234000000000000"},
		F120: &Lllnumeric{"Another test text"},
	}

	iso := Message{"0100", ASCII, true, data}

	res, err := iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	sample := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 54, 52, 51, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 54, 52, 51, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	if bytes.Compare(res, sample) != 0 {
		t.Error("ISO Encode error!")
	}
}

func TestDecode(t *testing.T) {

	input := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 54, 52, 51, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 54, 52, 51, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	// init empty iso message struct
	iso := Message{"", ASCII, true, newDataIso()}

	// parse data from bytes to iso struct
	err := iso.Load(input)

	if err != nil {
		t.Error("ISO Decode error:", err)
	}

	var res []byte

	// set second bitmap because field 120 in struct (need if more than 63 fields in message)
	iso.SecondBitmap = true

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

	sample := []byte{48, 49, 48, 48, 114, 60, 36, 129, 40, 224, 152, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 54, 52, 51, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 54, 52, 51, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}

	if bytes.Compare(res, sample) != 0 {
		t.Error("ISO Encode error!")
	}
}

// готовим контейнер для загрузки данных
func newDataIso() *DataIso {
	return &DataIso{
		F2 : &Llnumeric{},
		F3 : &Numeric{},
		F4 : &Numeric{},
		F7 : &Numeric{},
		F11: &Numeric{},
		F12: &Numeric{},
		F13: &Numeric{},
		F14: &Numeric{},
		F19: &Numeric{},
		F22: &Numeric{},
		F25: &Numeric{},
		F32: &Llnumeric{},
		F35: &Llnumeric{},
		F37: &Alphanumeric{},
		F39: &Alphanumeric{},
		F41: &Alphanumeric{},
		F42: &Alphanumeric{},
		F43: &Alphanumeric{},
		F49: &Alphanumeric{},
		F52: NewBinary(nil),
		F53: &Numeric{},
		F120: &Lllnumeric{},
	}
}