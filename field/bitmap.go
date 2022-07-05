package field

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
)

const minBitmapLength = 8 // 64 bit, 8 bytes, or 16 hex digits
const defaultBitmaps = 1

var _ Field = (*Bitmap)(nil)

// NOTE: Bitmap does not support JSON encoding or decoding.
type Bitmap struct {
	spec    *Spec
	bitmap  *utils.Bitmap
	data    *Bitmap
	mapSize int
}

func NewBitmap(spec *Spec) *Bitmap {
	return &Bitmap{
		spec:    spec,
		mapSize: defaultBitmaps,
		bitmap:  utils.NewBitmap(64 * defaultBitmaps),
	}
}

func (f *Bitmap) Spec() *Spec {
	return f.spec
}

func (f *Bitmap) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Bitmap) SetMapSize(size int) {
	f.mapSize = size
	f.Reset()
}

func (f *Bitmap) SetBytes(b []byte) error {
	f.bitmap = utils.NewBitmapFromData(b)
	if f.data != nil {
		*(f.data) = *f
	}
	return nil
}

func (f *Bitmap) Bytes() ([]byte, error) {
	return f.bitmap.Bytes(), nil
}

func (f *Bitmap) String() (string, error) {
	return f.bitmap.String(), nil
}

func (f *Bitmap) Pack() ([]byte, error) {
	f.setBitmapFields()

	count := f.bitmapsCount()

	// here we have max possible bytes for the bitmap 8*mapSize
	data, err := f.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bytes: %w", err)
	}

	data = data[0 : 8*count]

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	return packed, nil
}

// Unpack of the Bitmap field returns data of varied length
// if there is only primary bitmap (bit 1 is not set) we return only 8 bytes (or 16 for hex encoding)
// if secondary bitmap presents (bit 1 is set) we return 16 bytes (or 32 for hex encoding)
// and so on for mapSize
func (f *Bitmap) Unpack(data []byte) (int, error) {
	minLen, _, err := f.spec.Pref.DecodeLength(minBitmapLength, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}

	rawBitmap := make([]byte, 0)
	read := 0

	// read max
	for i := 0; i < f.getMapSize(); i++ {
		decoded, readDecoded, err := f.spec.Enc.Decode(data[read:], minLen)
		if err != nil {
			return 0, fmt.Errorf("failed to decode content for %d bitmap: %w", i+1, err)
		}
		read += readDecoded

		rawBitmap = append(rawBitmap, decoded...)
		bitmap := utils.NewBitmapFromData(decoded)

		// if no more bitmaps, exit loop
		if !bitmap.IsSet(1) {
			break
		}

		if i == f.getMapSize()-1 && bitmap.IsSet(1) {
			return 0, fmt.Errorf("failed to decode content for %d bitmap: invalid extended bitmap indicator", i+1)
		}
	}

	if err := f.SetBytes(rawBitmap); err != nil {
		return 0, fmt.Errorf("failed to set bytes: %w", err)
	}

	return read, nil
}

func (f *Bitmap) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	bmap, ok := v.(*Bitmap)
	if !ok {
		return fmt.Errorf("data does not match required *Bitmap type")
	}

	bmap.bitmap = f.bitmap

	return nil
}

func (f *Bitmap) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	bmap, ok := data.(*Bitmap)
	if !ok {
		return fmt.Errorf("data does not match required *Bitmap type")
	}

	f.data = bmap
	if bmap.bitmap != nil {
		f.bitmap = bmap.bitmap
	}
	return nil
}

func (f *Bitmap) Marshal(data interface{}) error {
	return f.SetData(data)
}

func (f *Bitmap) Reset() {
	f.bitmap = utils.NewBitmap(64 * f.getMapSize())
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

func (f *Bitmap) bitmapsCount() int {
	count := 1
	for i := 0; i < f.getMapSize(); i++ {
		if f.IsSet(i*64 + 1) {
			count += 1
		}
	}

	return count
}

func (f *Bitmap) getMapSize() int {
	if f.mapSize == 0 {
		f.mapSize = defaultBitmaps
		f.Reset()
	}
	return f.mapSize
}

func (f *Bitmap) setBitmapFields() bool {
	// 2nd bitmap bits 65 -128
	// bitmap bit 1

	// 3rd bitmap bits 129-192
	// bitmap bit 65

	// ...

	// start from the 2nd bitmap as for the 1st bitmap we don't need to set any bits
	for bitmapIndex := f.getMapSize(); bitmapIndex > 1; bitmapIndex-- {

		// are there fields for this (bitmapIndex) bitmap?
		bitmapStart := (bitmapIndex-1)*64 + 2 // we skip firt bit as it's for the next bitmap
		bitmapEnd := (bitmapIndex) * 64       //

		for i := bitmapStart; i <= bitmapEnd; i++ {
			bitmapBit := (bitmapIndex-2)*64 + 1
			if f.IsSet(i) {
				for subBitmap := bitmapBit; subBitmap > 0; subBitmap -= 64 {
					f.Set(subBitmap)
				}
				return true
			}
		}
	}

	return false
}
