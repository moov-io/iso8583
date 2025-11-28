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
// TODO: delete as we don't use it anymore
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

// Creates a map with new instances of Fields (Field interface)
// based on the field type in MessageSpec.
func (s *MessageSpec) CreateMessageFields() map[int]field.Field {

	fields := map[int]field.Field{}

	for k, specField := range s.Fields {
		fields[k] = field.NewInstanceOf(specField)
	}

	return fields
}
