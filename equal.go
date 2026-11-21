package iso8583

import (
	"fmt"
	"maps"
	"slices"

	"github.com/moov-io/iso8583/field"
)

// Equal reports whether messages a and b have the same MTI and the same set
// fields, with every pair of fields equal according to field.Equal.
//
// Two nil messages are equal. The primary bitmap is not compared because it
// is derived from the set fields. Messages built from different specs can be
// equal if their field values are equal.
func Equal(a, b *Message) (bool, error) {
	if a == nil || b == nil {
		return a == nil && b == nil, nil
	}
	if a == b {
		return true, nil
	}

	af := withoutBitmap(a.GetFields())
	bf := withoutBitmap(b.GetFields())
	if len(af) != len(bf) {
		return false, nil
	}

	for _, id := range slices.Sorted(maps.Keys(af)) {
		other, ok := bf[id]
		if !ok {
			return false, nil
		}
		eq, err := field.Equal(af[id], other)
		if err != nil {
			return false, fmt.Errorf("comparing field %d: %w", id, err)
		}
		if !eq {
			return false, nil
		}
	}

	return true, nil
}

func withoutBitmap(fields map[int]field.Field) map[int]field.Field {
	delete(fields, bitmapIdx)

	return fields
}
