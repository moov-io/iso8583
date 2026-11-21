package field

import (
	"bytes"
	"fmt"
	"maps"
	"reflect"
	"slices"
)

// Equal reports whether fields a and b hold the same value.
//
// Two nil fields are equal, and a nil field never equals a non-nil one.
// Composite fields are equal when the same subfields are set and every pair
// of subfields is Equal; the composite bitmap is not compared because it is
// derived from the set subfields. All other fields are compared by their
// value (Bytes), which ignores the length prefix and encoding in the spec.
//
// Equal compares values, not types: a String and a Numeric that hold "123"
// are equal. It is intended for tests and for comparing messages built from
// the same spec.
func Equal(a, b Field) (bool, error) {
	if isNil(a) || isNil(b) {
		return isNil(a) && isNil(b), nil
	}

	ca, aIsComposite := a.(*Composite)
	cb, bIsComposite := b.(*Composite)
	if aIsComposite != bIsComposite {
		return false, nil
	}
	if aIsComposite {
		return compositeEqual(ca, cb)
	}

	av, err := a.Bytes()
	if err != nil {
		return false, fmt.Errorf("getting bytes of first field: %w", err)
	}
	bv, err := b.Bytes()
	if err != nil {
		return false, fmt.Errorf("getting bytes of second field: %w", err)
	}

	return bytes.Equal(av, bv), nil
}

func compositeEqual(a, b *Composite) (bool, error) {
	if a == b {
		return true, nil
	}

	as := withoutBitmap(a.GetSubfields())
	bs := withoutBitmap(b.GetSubfields())
	if len(as) != len(bs) {
		return false, nil
	}

	for _, id := range slices.Sorted(maps.Keys(as)) {
		bf, ok := bs[id]
		if !ok {
			return false, nil
		}
		eq, err := Equal(as[id], bf)
		if err != nil {
			return false, fmt.Errorf("comparing subfield %s: %w", id, err)
		}
		if !eq {
			return false, nil
		}
	}

	return true, nil
}

func withoutBitmap(fields map[string]Field) map[string]Field {
	for id, f := range fields {
		if _, ok := f.(*Bitmap); ok {
			delete(fields, id)
		}
	}

	return fields
}

func isNil(f Field) bool {
	if f == nil {
		return true
	}
	v := reflect.ValueOf(f)

	return v.Kind() == reflect.Ptr && v.IsNil()
}
