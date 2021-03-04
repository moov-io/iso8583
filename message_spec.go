package iso8583

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"

	"github.com/moov-io/iso8583/field"
	"github.com/stoewer/go-strcase"
)

type MessageSpec struct {
	Fields map[int]field.Field
}

// Creates a map with new instances of Fields (Field interface)
// based on the field type in MessageSpec.
func (s *MessageSpec) CreateMessageFields() map[int]field.Field {

	fields := map[int]field.Field{}

	for k, specField := range s.Fields {
		fields[k] = s.createMessageField(specField)
	}

	return fields
}

// Get field's index by identifier
func (s *MessageSpec) GetFieldIndex(identifier, stringCase string) (int, error) {

	for key, f := range s.Fields {

		newIdentifier := strcase.SnakeCase(f.Spec().GetIdentifier())
		if stringCase == "CamelCase" {
			newIdentifier = strcase.UpperCamelCase(f.Spec().GetIdentifier())
		}

		if newIdentifier == identifier {
			return key, nil
		}
	}

	return 0, errors.New("don't find any specification by identifier")
}

// Customize unmarshal of json
func (s *MessageSpec) UnmarshalJSON(b []byte) error {
	dummy := newMessageSpecMarshaller()

	err := json.Unmarshal(b, &dummy.Fields)
	if err != nil {
		return errors.New("failed to parse json string of message")
	}

	for key, element := range dummy.Fields {
		var index int
		_, err := fmt.Sscanf(key, "%03d", &index)
		if err != nil {
			fmt.Println("There is invalid field of specification:", key)
			continue
		}

		err = s.createMessageSpecFieldWithMarshaller(index, element)
		if err != nil {
			fmt.Println("failed to create specification field:", err)
			continue
		}
	}
	return nil
}

// Customize marshal of json
func (s *MessageSpec) MarshalJSON() ([]byte, error) {
	dummy, err := s.createMarshaller()
	if err != nil {
		return nil, err
	}
	return json.Marshal(dummy.Fields)
}

// Customize marshal of xml
func (s *MessageSpec) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

	dummy := newMessageSpecMarshaller()

	for {
		var element specMarshaller
		err := decoder.Decode(&element)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		dummy.FieldArray = append(dummy.FieldArray, element)
	}

	for _, element := range dummy.FieldArray {
		var index int
		_, err := fmt.Sscanf(element.XMLName.Local, "F%03d", &index)
		if err != nil {
			fmt.Println("There is invalid field of specification:", element.XMLName.Local)
			continue
		}

		err = s.createMessageSpecFieldWithMarshaller(index, element)
		if err != nil {
			fmt.Println("failed to create specification field:", err)
			continue
		}
	}

	return nil
}

// Customize unmarshal of xml
func (s *MessageSpec) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	dummy, err := s.createMarshaller()
	if err != nil {
		return err
	}
	start.Name = dummy.XMLName
	return encoder.EncodeElement(dummy, start)
}

func (s *MessageSpec) createMessageSpecFieldWithMarshaller(index int, marshaller specMarshaller) (err error) {
	newField := marshaller.createSpecificationField()

	if newField == nil {
		return errors.New("failed to create specification field")
	}

	if s.Fields == nil {
		s.Fields = make(map[int]field.Field, 0)
	}
	s.Fields[index] = marshaller.createSpecificationField()

	return nil
}

func (s *MessageSpec) createMarshaller() (*messageSpecMarshaller, error) {
	dummy := newMessageSpecMarshaller()

	var fieldIndexes []int
	for id := range s.Fields {
		fieldIndexes = append(fieldIndexes, id)
	}

	sort.Ints(fieldIndexes)
	for _, index := range fieldIndexes {
		field := s.Fields[index]

		if err := dummy.addSpecMarshaller(index, field); err != nil {
			return nil, err
		}
	}

	return dummy, nil
}

func (s *MessageSpec) createMessageField(specField field.Field) field.Field {
	fieldType := reflect.TypeOf(specField).Elem()

	// create new field and convert it to field.Field interface
	fl := reflect.New(fieldType).Interface().(field.Field)
	fl.SetSpec(specField.Spec())

	return fl
}
