package field

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*Bitmap)(nil)

type Bitmap struct {
	spec   *Spec
	bitmap *utils.Bitmap
}

func NewBitmap(spec *Spec) Field {
	return &Bitmap{
		spec:   spec,
		bitmap: utils.NewBitmap(192),
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
	if f.isSecondary() {
		f.bitmap.Set(1)
	}

	data := f.Bytes()

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %v", err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %v", err)
	}

	if !f.bitmap.IsSet(1) {
		packed = packed[:len(packed)/2]
	}

	return append(packedLength, packed...), nil
}

// Unpack of the Bitmap field returns data of varied length
// if there is only primary bitmap (bit 1 is not set) we return only 8 bytes
// if secondary bitmap presents (bit 1 is set) we return 16 bytes

const minBitmapLength = 8 // 64 bit, 8 bytes, or 16 hex digits
const maxBitmaps = 3

func (f *Bitmap) Unpack(data []byte) (int, error) {
	minLen, err := f.spec.Pref.DecodeLength(minBitmapLength, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %v", err)
	}

	rawBitmap := make([]byte, 0)
	read := 0

	// read max
	for i := 0; i < maxBitmaps; i++ {
		start := i * minLen
		end := (i + 1) * minLen

		if len(data) < end {
			return 0, fmt.Errorf("not enough data to read %d bitmap", i+1)
		}

		decoded, err := f.spec.Enc.Decode(data[start:end], 0)
		if err != nil {
			return 0, fmt.Errorf("failed to decode content: %v", err)
		}
		read += minLen

		rawBitmap = append(rawBitmap, decoded...)
		bitmap := utils.NewBitmapFromData(decoded)

		// if no more bitmaps, exit loop
		if !bitmap.IsSet(1) {
			break
		}
	}

	f.bitmap = utils.NewBitmapFromData(rawBitmap)
	return read, nil
}

func (f *Bitmap) Reset() {
	f.bitmap = utils.NewBitmap(192)
}

func (f *Bitmap) Set(i int) {
	f.bitmap.Set(i)
}

func (f *Bitmap) IsSet(i int) bool {
	return f.bitmap.IsSet(i)
}

func (f *Bitmap) Len() int {
	return f.bitmap.Len()
}

func (f *Bitmap) isSecondary() bool {
	for i := 65; i <= 128; i++ {
		if f.IsSet(i) {
			return true
		}
	}

	return false
}
