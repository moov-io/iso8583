package field

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*Bitmap)(nil)

type Bitmap struct {
	Value  []byte
	spec   *Spec
	bitmap *utils.Bitmap
	maxId  int
}

func NewBitmap(spec *Spec) Field {
	return &Bitmap{
		spec:   spec,
		bitmap: utils.NewBitmap(128),
	}
}

func (f *Bitmap) Spec() *Spec {
	return f.spec
}

func (f *Bitmap) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Bitmap) SetBytes(b []byte) {
	f.bitmap = utils.NewBitmapFromData(b)
}

func (f *Bitmap) Bytes() []byte {
	return f.bitmap.Bytes()
}

func (f *Bitmap) String() string {
	return f.bitmap.String()
}

func (f *Bitmap) Pack() ([]byte, error) {
	fmt.Println("pack bitmap")

	// let's test and read more about it
	if f.maxId > 64 {
		f.bitmap.Set(1)
	}

	data := f.Bytes()

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", f.spec.Description, err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", f.spec.Description, err)
	}

	if !f.bitmap.IsSet(1) {
		packed = packed[:len(packed)/2]
	}

	return append(packedLength, packed...), nil
}

// Unpack of the Bitmap field returns data of varied length
// if there is only primary bitmap (bit 1 is not set) we return only 8 bytes
// if secondary bitmap presents (bit 1 is set) we return 16 bytes
func (f *Bitmap) Unpack(data []byte) ([]byte, int, error) {
	dataLen, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	start := f.spec.Pref.Length()
	end := f.spec.Pref.Length() + dataLen
	raw, err := f.spec.Enc.Decode(data[start:end], 0)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", f.spec.Description, err)
	}

	bitmap := utils.NewBitmapFromData(raw)

	if bitmap.IsSet(1) {
		f.bitmap = utils.NewBitmapFromData(raw[:16])
		return raw[:16], dataLen, nil
	}

	f.bitmap = utils.NewBitmapFromData(raw[:8])
	return raw[:8], dataLen / 2, nil
}

func (f *Bitmap) Reset() {
	f.bitmap = utils.NewBitmap(128)
}

func (f *Bitmap) Set(i int) {
	if i > f.maxId {
		f.maxId = i
	}

	f.bitmap.Set(i)
}

func (f *Bitmap) IsSet(i int) bool {
	return f.bitmap.IsSet(i)
}

func (f *Bitmap) Len() int {
	return f.bitmap.Len()
}
