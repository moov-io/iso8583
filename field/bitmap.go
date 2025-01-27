package field

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var _ Field = (*Bitmap)(nil)

// Bitmap is a 1-indexed big endian bitmap field.
type Bitmap struct {
	spec         *Spec
	data         []byte
	bitmapLength int
}

const defaultBitmapLength = 8

const firstBitOn = 0b10000000 // big endian

func NewBitmap(spec *Spec) *Bitmap {
	length := spec.Length
	if length == 0 {
		length = defaultBitmapLength
	}

	return &Bitmap{
		spec:         spec,
		data:         make([]byte, length),
		bitmapLength: length,
	}
}

func (f *Bitmap) Spec() *Spec {
	return f.spec
}

func (f *Bitmap) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Bitmap) SetBytes(b []byte) error {
	f.data = b
	return nil
}

func (f *Bitmap) Bytes() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return f.data, nil
}

func (f *Bitmap) String() (string, error) {
	if f == nil {
		return "", nil
	}

	var bits []string

	for _, byte_ := range f.data {
		bits = append(bits, fmt.Sprintf("%08b", byte_))
	}

	return strings.Join(bits, " "), nil
}

func (f *Bitmap) Pack() ([]byte, error) {
	packed, err := f.spec.Enc.Encode(f.data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	return packed, nil
}

// Unpack sets the bitmap data. It returns the number of bytes read from the
// data. Usually it's 8 for binary, 16 for hex - for a single bitmap.
// If DisableAutoExpand is not set (default), it will read all bitmaps until
// the first bit of the read bitmap is not set.
// If DisableAutoExpand is set, it will only read the first bitmap regardless
// of the first bit being set.
func (f *Bitmap) Unpack(data []byte) (int, error) {
	minLen, _, err := f.spec.Pref.DecodeLength(f.bitmapLength, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}

	f.data = make([]byte, 0)
	read := 0

	var i int

	// read until we have no more bitmaps
	// or only read one bitmap if DisableAutoExpand is set
	for {
		i++
		decoded, readDecoded, err := f.spec.Enc.Decode(data[read:], minLen)
		if err != nil {
			return 0, fmt.Errorf("failed to decode content for %d bitmap: %w", i, err)
		}
		read += readDecoded
		f.data = append(f.data, decoded...)

		// if it's a fixed bitmap or first bit of the decoded bitmap is not set, exit loop
		if f.spec.DisableAutoExpand || decoded[0]&firstBitOn == 0 {
			break
		}
	}

	return read, nil
}

// Deprecated. Use Marshal instead
func (f *Bitmap) SetData(data interface{}) error {
	return f.Marshal(data)
}

func (f *Bitmap) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	bmap, ok := v.(*Bitmap)
	if !ok {
		return fmt.Errorf("unsupported type: expected *Bitmap, got %T", v)
	}

	bmap.data = f.data

	return nil
}

func (f *Bitmap) Marshal(v interface{}) error {
	if v == nil {
		return nil
	}

	bmap, ok := v.(*Bitmap)
	if !ok {
		return fmt.Errorf("data does not match required *Bitmap type")
	}

	f.data = bmap.data
	return nil
}

// Reset the bitmap to its initial state because of how message works,
// Message need a way to initialize bitmap. That's why we set parameters to
// their default values here like we do in constructor.
func (f *Bitmap) Reset() {
	length := f.spec.Length
	if length == 0 {
		length = defaultBitmapLength
	}

	f.bitmapLength = length
	// this actually resets the bitmap
	f.data = make([]byte, f.bitmapLength)
}

// For auto expand mode if we expand bitmap we should set bit that shows the presence of the next bitmap
func (f *Bitmap) Set(n int) {
	if n <= 0 {
		return
	}

	// do we have to expand bitmap?
	if n > len(f.data)*8 {
		if f.spec.DisableAutoExpand {
			return
		}

		// calculate how many bitmaps we need to store n-th bit
		bitmapIndex := (n - 1) / (f.bitmapLength * 8)
		newBitmapsCount := (bitmapIndex + 1)

		// set first bit of the first byte of the last bitmap in
		// current data to 1 to show the presence of the next bitmap
		f.data[len(f.data)-f.bitmapLength] |= firstBitOn

		// add new empty bitmaps and for every new bitmap except the
		// last one, set bit that shows the presence of the next bitmap
		for i := newBitmapsCount - len(f.data)/f.bitmapLength; i > 0; i-- {
			newBitmap := make([]byte, f.bitmapLength)
			// set first bit of the first byte of the new bitmap to 1
			// but only if it is not the last bitmap
			if i > 1 {
				newBitmap[0] = firstBitOn
			}
			f.data = append(f.data, newBitmap...)
		}
	}

	// set bit
	f.data[(n-1)/8] |= 1 << (uint(7-(n-1)) % 8)
}

func (f *Bitmap) IsSet(n int) bool {
	if n <= 0 || n > len(f.data)*8 {
		return false
	}

	return f.data[(n-1)/8]&(1<<(uint(7-(n-1))%8)) != 0
}

func (f *Bitmap) Len() int {
	return len(f.data) * 8
}

func (f *Bitmap) MarshalJSON() ([]byte, error) {
	data, err := f.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bytes: %v", err)
	}
	return json.Marshal(strings.ToUpper(hex.EncodeToString(data)))
}

// Takes in a HEX based string
func (f *Bitmap) UnmarshalJSON(b []byte) error {
	unqouted, err := strconv.Unquote(string(b))
	if err != nil {
		return fmt.Errorf("failed to unquote input: %w", err)
	}

	bs, err := hex.DecodeString(unqouted)
	if err != nil {
		return fmt.Errorf("failed to decode hex string: %w", err)
	}

	return f.SetBytes(bs)
}

// IsBitmapPresenceBit checks if the bit at position n in the bitmap is an
// indicator of the presence of an additional bitmap. For fixed-length bitmaps
// (when DisableAutoExpand is set in the specification), this method will
// always return false since additional bitmaps are not applicable.
func (f *Bitmap) IsBitmapPresenceBit(n int) bool {
	// there are not presence bits in fixed bitmaps
	if f.spec.DisableAutoExpand {
		return false
	}

	if n <= 0 {
		return false
	}

	// check if n is the first bit of a bitmap
	if n%(f.bitmapLength*8) == 1 {
		return true
	}

	return false
}
