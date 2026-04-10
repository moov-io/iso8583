package iso8583

import (
	"fmt"
	"strconv"

	iso8583errors "github.com/moov-io/iso8583/errors"
	"github.com/moov-io/iso8583/field"
)

// MessageScanner provides forward-only cursor-based parsing of an ISO 8583
// message. It is designed for cases where only a few fields need to be
// inspected (e.g., MTI and STAN in a proxy) without unpacking the entire
// message.
//
// MTI and bitmap are parsed automatically as needed. Fields can only be
// scanned in ascending order; attempting to scan a field at or before the
// current cursor position returns an error.
type MessageScanner struct {
	spec   *MessageSpec
	src    []byte
	offset int

	bitmap        *field.Bitmap
	lastID        int // last scanned field ID, -1 initially
	scannedBitmap bool
}

// NewMessageScanner creates a new scanner for the given spec and raw message
// bytes. It does not perform any parsing until ScanField is called.
func NewMessageScanner(spec *MessageSpec, src []byte) *MessageScanner {
	return &MessageScanner{
		spec:   spec,
		src:    src,
		lastID: -1,
	}
}

// ScanField parses the message up to and including the field with the given
// id, returning the parsed field. Fields between the current cursor position
// and the target field are consumed but discarded.
//
// Field 0 (MTI) must be scanned before any other field. Fields must be
// scanned in strictly ascending order. Any error is returned as
// *iso8583errors.UnpackError.
func (s *MessageScanner) ScanField(id int) (field.Field, error) {
	f, err := s.scanField(id)
	if err != nil {
		return nil, &iso8583errors.UnpackError{
			Err:        err,
			FieldID:    strconv.Itoa(id),
			RawMessage: s.src,
		}
	}
	return f, nil
}

func (s *MessageScanner) scanField(id int) (field.Field, error) {
	if id < 0 {
		return nil, fmt.Errorf("invalid field id: %d", id)
	}

	if id <= s.lastID {
		return nil, fmt.Errorf("field %d is at or before current position %d: scanner is forward-only", id, s.lastID)
	}

	// Ensure MTI is scanned; return it if that's what was requested.
	if s.lastID == -1 {
		f, err := s.scanMTI()
		if err != nil {
			return nil, fmt.Errorf("failed to scan MTI: %w", err)
		}
		if id == mtiIdx {
			return f, nil
		}
	}

	// Ensure bitmap is scanned; return it if that's what was requested.
	if !s.scannedBitmap {
		if err := s.parseBitmap(); err != nil {
			return nil, fmt.Errorf("failed to parse bitmap: %w", err)
		}
		if id == bitmapIdx {
			return s.bitmap, nil
		}
	}

	for i := s.nextDataField(); i <= id && i <= s.bitmap.Len(); i = s.nextDataField() {
		if s.bitmap.IsBitmapPresenceBit(i) {
			s.lastID = i
			continue
		}

		if !s.bitmap.IsSet(i) {
			s.lastID = i
			continue
		}

		specField, ok := s.spec.Fields[i]
		if !ok {
			return nil, fmt.Errorf("field %d is not defined in the spec", i)
		}

		f := field.NewInstanceOf(specField)

		read, err := f.Unpack(s.src[s.offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to unpack field %d (%s): %w", i, f.Spec().Description, err)
		}

		s.offset += read
		s.lastID = i

		if i == id {
			return f, nil
		}
	}

	// Field was not present in the bitmap
	return nil, fmt.Errorf("field %d is not set in the bitmap", id)
}

func (s *MessageScanner) scanMTI() (field.Field, error) {
	specField, ok := s.spec.Fields[mtiIdx]
	if !ok {
		return nil, fmt.Errorf("field %d (MTI) is not defined in the spec", mtiIdx)
	}

	f := field.NewInstanceOf(specField)

	read, err := f.Unpack(s.src[s.offset:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack MTI: %w", err)
	}

	s.offset += read
	s.lastID = mtiIdx

	return f, nil
}

func (s *MessageScanner) parseBitmap() error {
	specField, ok := s.spec.Fields[bitmapIdx]
	if !ok {
		return fmt.Errorf("field %d (bitmap) is not defined in the spec", bitmapIdx)
	}

	bmp, ok := field.NewInstanceOf(specField).(*field.Bitmap)
	if !ok {
		return fmt.Errorf("field %d is not a bitmap", bitmapIdx)
	}

	bmp.Reset()

	read, err := bmp.Unpack(s.src[s.offset:])
	if err != nil {
		return fmt.Errorf("failed to unpack bitmap: %w", err)
	}

	s.offset += read
	s.bitmap = bmp
	s.scannedBitmap = true
	s.lastID = bitmapIdx

	return nil
}

// nextDataField returns the next field ID to process after lastID.
func (s *MessageScanner) nextDataField() int {
	return max(s.lastID+1, 2)
}
