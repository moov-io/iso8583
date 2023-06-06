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

// FieldContainer should be implemented by by the type to be described
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
			fmt.Fprintf(w, fmt.Sprintf("F%-3s %s SUBFIELDS:\n", i, desc))
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

		fmt.Fprintf(w, fmt.Sprintf("F%-3s %s\t: %%s\n", i, desc), str)
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
	numericKeys := make([]int, 0)
	nonNumericKeys := make([]string, 0)

	for k := range fields {
		if id, err := strconv.Atoi(k); err == nil {
			numericKeys = append(numericKeys, id)
		} else {
			nonNumericKeys = append(nonNumericKeys, k)
		}
	}

	// Sorting numeric and non-numeric keys separately
	sort.Ints(numericKeys)
	sort.Strings(nonNumericKeys)

	keys := make([]string, 0, len(fields))

	// Appending numeric keys first (as strings)
	for _, key := range numericKeys {
		keys = append(keys, strconv.Itoa(key))
	}

	// Appending non-numeric keys
	keys = append(keys, nonNumericKeys...)

	return keys
}

// splitAndAnnotate splits bits blocks and annotates them with bit numbers
// and splits them by 32 bits if needed
func splitAndAnnotate(bits string) string {
	bitBlocks := strings.Split(bits, " ")
	if len(bitBlocks) == 0 {
		return ""
	}

	annotatedBits := make([]string, len(bitBlocks))
	bitsCount := len(bitBlocks[0])

	for i, block := range bitBlocks {
		startBit := i*bitsCount + 1
		endBit := (i + 1) * bitsCount
		annotatedBits[i] = fmt.Sprintf("[%d-%d]%s", startBit, endBit, block)

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
