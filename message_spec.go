package iso8583

import (
	"fmt"

	"github.com/moov-io/iso8583/field"
)

type MessageSpec struct {
	Name   string
	Fields map[int]field.Field
}

// Validate checks if the MessageSpec is valid.
func (s *MessageSpec) Validate() error {
	// we require MTI and Bitmap fields
	if _, ok := s.Fields[mtiIdx]; !ok {
		return fmt.Errorf("MTI field (%d) is required", mtiIdx)
	}

	if _, ok := s.Fields[bitmapIdx]; !ok {
		return fmt.Errorf("Bitmap field (%d) is required", bitmapIdx)
	}

	// check type of the bitmap field
	if _, ok := s.Fields[bitmapIdx].(*field.Bitmap); !ok {
		return fmt.Errorf("Bitmap field (%d) must be of type *field.Bitmap", bitmapIdx)
	}

	return nil
}
