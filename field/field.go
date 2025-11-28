package field

import "reflect"

// PathMarshaler provides the ability to marshal field values using path notation.
// The path uses dot notation (e.g., "11.1" or "3.2.1") to navigate nested
// composite fields and marshal values at any depth within the field hierarchy.
type PathMarshaler interface {
	MarshalPath(path string, v any) error
}

// PathUnmarshaler provides the ability to unmarshal field values using path
// notation. The path uses dot notation (e.g., "11.1" or "3.2.1") to navigate
// nested composite fields and unmarshal values from any depth within the field
// hierarchy.
type PathUnmarshaler interface {
	UnmarshalPath(path string, v any) error
}

// PathUnsetter provides the ability to unset fields using path notation.
// The path uses dot notation (e.g., "11.1" or "3.2.1") to navigate nested
// composite fields and unset values at any depth. Unset fields are replaced
// with zero-valued fields and excluded from operations like Pack() or Marshal().
type PathUnsetter interface {
	UnsetPath(idPaths ...string) error
}

type Field interface {
	// Spec returns the field spec
	Spec() *Spec

	// SetSpec sets the field spec
	SetSpec(spec *Spec)

	// Pack serializes field value into binary representation according
	// to the field spec
	Pack() ([]byte, error)

	// Unpack deserialises the field by reading length prefix and reading
	// corresponding number of bytes from the provided data parameter and
	// then decoding it according to the field spec
	Unpack(data []byte) (int, error)

	// SetBytes sets the field Value using its binary representation
	// provided in the data parameter
	SetBytes(data []byte) error

	// Bytes returns binary representation of the field Value
	Bytes() ([]byte, error)

	// Deprecated. Use Marshal instead.
	SetData(data interface{}) error

	// Unmarshal sets field Value into provided v. If v is nil or not
	// a pointer it returns error.
	Unmarshal(v interface{}) error

	// Marshal sets field Value from provided v. If v is nil or not
	// a pointer it returns error.
	Marshal(v interface{}) error

	// String returns a string representation of the field Value
	String() (string, error)
}

// NewInstanceOf creates a new instance of the same type as specField
// and sets its spec to be the same as specField's spec.
func NewInstanceOf(specField Field) Field {
	fieldType := reflect.TypeOf(specField).Elem()

	// create new field and convert it to field.Field interface

	//nolint:forcetypeassert // we know the type of the field we're creating here
	fl := reflect.New(fieldType).Interface().(Field)
	fl.SetSpec(specField.Spec())

	// if it's a composite field, we have to recusively create its subfields as well
	// TODO: remove this, as fields will be lazily constructed when accessed
	if composite, ok := fl.(CompositeWithSubfields); ok {
		composite.ConstructSubfields()
	}

	return fl
}
