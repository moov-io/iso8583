package iso8583

import (
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/moov-io/iso8583/field"
)

var defaultSpecName = "ISO 8583"

// FieldContainer should be implemented by the type to be described
// we use GetSubfields() as a common method to get subfields
// while Message doesn't implement FieldContainer interface directly
// we use MessageWrapper to wrap Message and implement FieldContainer
type FieldContainer interface {
	GetSubfields() map[string]field.Field
}

type ContainerWithBitmap interface {
	Bitmap() *field.Bitmap
}

// MessageWrapper implements FieldContainer interface for the iso8583.Message
// as currently it has GetFields() and not GetSubfields and it returns
// map[int]field.Field (key is int, not string)
type MessageWrapper struct {
	*Message
}

func (m *MessageWrapper) GetSubfields() map[string]field.Field {
	fields := m.Message.GetFields()
	result := make(map[string]field.Field, len(fields))
	for k, v := range fields {
		result[fmt.Sprintf("%d", k)] = v
	}
	return result
}

func Describe(message *Message, w io.Writer, filters ...FieldFilter) error {
	specName := defaultSpecName
	if spec := message.GetSpec(); spec != nil && spec.Name != "" {
		specName = spec.Name
	}
	fmt.Fprintf(w, "%s Message:\n", specName)

	tw := tabwriter.NewWriter(w, 0, 0, 2, '.', 0)

	mti, err := message.GetMTI()
	if err != nil {
		return fmt.Errorf("getting MTI: %w", err)
	}
	fmt.Fprintf(tw, "MTI\t: %s\n", mti)

	// use default filter
	if len(filters) == 0 {
		filters = DefaultFilters()
	}

	err = DescribeFieldContainer(&MessageWrapper{message}, tw, filters...)
	if err != nil {
		return fmt.Errorf("describing message: %w", err)
	}

	tw.Flush()

	return nil
}

// DescribeFieldContainer describes the FieldContainer (e.g. Wrapped Message or CompositeField)
func DescribeFieldContainer(container FieldContainer, w io.Writer, filters ...FieldFilter) error {
	// making filter map
	filterMap := make(map[string]FilterFunc)

	for _, filter := range filters {
		filter(filterMap)
	}

	var errorList []string

	// container may have bitmap
	var bitmap *field.Bitmap
	if container, ok := container.(ContainerWithBitmap); ok {
		bitmap = container.Bitmap()
	}

	if bitmap != nil {
		bitmapRaw, err := bitmap.Bytes()
		if err != nil {
			return fmt.Errorf("getting bitmap bytes: %w", err)
		}
		fmt.Fprintf(w, "Bitmap HEX\t: %s\n", strings.ToUpper(hex.EncodeToString(bitmapRaw)))

		bits, err := bitmap.String()
		if err != nil {
			return fmt.Errorf("getting bitmap: %w", err)
		}
		fmt.Fprintf(w, "Bitmap bits\t:\n%s\n", splitAndAnnotate(bits))
	}

	fields := container.GetSubfields()

	for _, i := range sortFieldIDs(fields) {
		f := fields[i]

		// skip bitmap as it's already displayed
		if f == bitmap {
			continue
		}

		desc := f.Spec().Description

		// check if field has subfields (e.g. CompositeField)
		if container, ok := f.(FieldContainer); ok {
			fmt.Fprintf(w, "DE%-3s %s SUBFIELDS:\n", i, desc)
			fmt.Fprintln(w, "-------------------------------------------")
			DescribeFieldContainer(container, w, filters...)
			fmt.Fprintln(w, "------------------------------------------")
			continue
		}

		// otherwise, print the field as usual

		str, err := f.String()
		if err != nil {
			errorList = append(errorList, err.Error())
			continue
		}

		// apply filtering
		if filter, existed := filterMap[i]; existed {
			str = filter(str, fields[i])
		}

		fmt.Fprintf(w, "DE%-3s %s\t: %s\n", i, desc, str)
	}

	if len(errorList) > 0 {
		fmt.Fprintf(w, "\nUnpacking Errors:\n")
		for _, err := range errorList {
			fmt.Fprintf(w, "- %s:\n", err)
		}

		return fmt.Errorf("displaying fields: %s", strings.Join(errorList, ","))
	}

	return nil
}

func sortFieldIDs(fields map[string]field.Field) []string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		// check if both keys are numeric
		ni, ei := strconv.Atoi(keys[i])
		nj, ej := strconv.Atoi(keys[j])
		if ei == nil && ej == nil {
			// if both keys are numeric, compare as integers
			return ni < nj
		}
		if ei == nil {
			// if only i is numeric, it goes first
			return true
		}
		if ej == nil {
			// if only j is numeric, it goes first
			return false
		}
		// if neither key is numeric, compare as strings
		return keys[i] < keys[j]
	})

	return keys
}

// splitAndAnnotate splits bits blocks and annotates them with bit numbers
// and splits them by 32 bits if needed
func splitAndAnnotate(bits string) string {
	if bits == "" {
		return ""
	}

	bitBlocks := strings.Split(bits, " ")

	annotatedBits := make([]string, len(bitBlocks))
	bitsCount := len(bitBlocks[0])

	pad := 0
	if len(bitBlocks) > 4 { // if multiple rows
		pad = 9 // pad to vertical align byte blocks
	}

	for i, block := range bitBlocks {
		startBit := i*bitsCount + 1
		endBit := (i + 1) * bitsCount
		pos := fmt.Sprintf("[%d-%d]", startBit, endBit)
		annotatedBits[i] = fmt.Sprintf("%*s%s", pad, pos, block)
		// split by 32 bits and check if it's not the last block
		isLastBlock := i == len(bitBlocks)-1
		isEndOf32Bits := endBit%32 == 0

		if isEndOf32Bits && !isLastBlock {
			annotatedBits[i] += "\n"
		} else if !isLastBlock {
			annotatedBits[i] += " "
		}
	}

	return strings.Join(annotatedBits, "")
}
