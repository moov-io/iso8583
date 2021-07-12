package describe

import (
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/field"
)

func Message(w io.Writer, message *iso8583.Message) error {
	fmt.Fprintf(w, "ISO 8583 Message:\n")

	printer := fieldPrinter{}

	mti, err := message.GetMTI()
	if err != nil {
		return fmt.Errorf("getting MTI: %w", err)
	}
	printer.addField("MTI", mti)

	bitmapRaw, err := message.Bitmap().Bytes()
	if err != nil {
		return fmt.Errorf("getting bitmap bytes: %w", err)
	}
	printer.addField("Bitmap", strings.ToUpper(hex.EncodeToString(bitmapRaw)))

	bitmap, err := message.Bitmap().String()
	if err != nil {
		return fmt.Errorf("getting bitmap: %w", err)
	}
	printer.addField("Bitmap bits", bitmap)

	// display the rest of all set fields
	fields := message.GetFields()
	for _, i := range sortFieldIDs(fields) {
		// skip the bitmap
		if i == 1 {
			continue
		}
		field := fields[i]
		desc := field.Spec().Description
		str, err := field.String()
		if err != nil {
			return fmt.Errorf("getting string value of field %d: %w", i, err)
		}
		printer.addField(fmt.Sprintf("F%03d %s", i, desc), str)
	}

	printer.print()

	return nil
}

func sortFieldIDs(fields map[int]field.Field) []int {
	keys := make([]int, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	return keys
}

type printableField struct {
	description string
	value       string
}
type fieldPrinter struct {
	fields []printableField
	maxLen int
}

func (p *fieldPrinter) addField(description, value string) {
	field := printableField{description, value}
	if descLen := utf8.RuneCountInString(description); descLen > p.maxLen {
		p.maxLen = descLen
	}

	p.fields = append(p.fields, field)
}

func (p *fieldPrinter) print() {
	// let's add some space after the description
	maxDescriptionLength := p.maxLen + 3
	for _, field := range p.fields {
		padding := strings.Repeat(".", maxDescriptionLength-utf8.RuneCountInString(field.description))
		fmt.Printf("%s%s: %v\n", field.description, padding, field.value)
	}
}
