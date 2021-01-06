package field

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*BitmapField)(nil)

type BitmapField struct {
	Value string
	spec  *Spec
}

func NewBitmapField(spec *Spec) Field {
	return &BitmapField{
		spec: spec,
	}
}

func (f *BitmapField) Spec() *Spec {
	return f.spec
}

func (f *BitmapField) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *BitmapField) SetBytes(b []byte) {
	f.Value = string(b)
}

func (f *BitmapField) Bytes() []byte {
	return []byte(f.Value)
}

func (f *BitmapField) String() string {
	return f.Value
}

func (f *BitmapField) Pack(data []byte) ([]byte, error) {
	bitmap := utils.NewBitmapFromData(data)

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", f.spec.Description, err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", f.spec.Description, err)
	}

	if !bitmap.IsSet(1) {
		packed = packed[:len(packed)/2]
	}

	return append(packedLength, packed...), nil
}

// Unpack of the Bitmap field returns data of varied length
// if there is only primary bitmap (bit 1 is not set) we return only 8 bytes
// if secondary bitmap presents (bit 1 is set) we return 16 bytes
func (f *BitmapField) Unpack(data []byte) ([]byte, int, error) {
	dataLen, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	start := f.spec.Pref.Length()
	end := f.spec.Pref.Length() + dataLen
	raw, err := f.spec.Enc.Decode(data[start:end])
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	bitmap := utils.NewBitmapFromData(raw)

	if bitmap.IsSet(1) {
		return raw[:16], dataLen, nil
	}

	return raw[:8], dataLen / 2, nil
}
