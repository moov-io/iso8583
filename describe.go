package iso8583

import (
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/moov-io/iso8583/field"
)

var defaultSpecName = "ISO 8583"

func Describe(message *Message, w io.Writer) error {
	specName := defaultSpecName
	if spec := message.GetSpec(); spec != nil && spec.Name != "" {
		specName = spec.Name
	}
	fmt.Fprintf(w, "%s Message:\n", specName)

	tw := tabwriter.NewWriter(w, 0, 0, 3, '.', 0)

	mti, err := message.GetMTI()
	if err != nil {
		return fmt.Errorf("getting MTI: %w", err)
	}
	fmt.Fprintf(tw, "MTI\t: %s\n", mti)

	bitmapRaw, err := message.Bitmap().Bytes()
	if err != nil {
		return fmt.Errorf("getting bitmap bytes: %w", err)
	}
	fmt.Fprintf(tw, "Bitmap\t: %s\n", strings.ToUpper(hex.EncodeToString(bitmapRaw)))

	bitmap, err := message.Bitmap().String()
	if err != nil {
		return fmt.Errorf("getting bitmap: %w", err)
	}
	fmt.Fprintf(tw, "Bitmap bits\t: %s\n", bitmap)

	// display the rest of all set fields
	fields := message.GetFields()

	var errorList []string

	for _, i := range sortFieldIDs(fields) {
		// skip the bitmap
		if i == 1 {
			continue
		}
		field := fields[i]
		desc := field.Spec().Description
		str, err := field.String()
		if err != nil {
			errorList = append(errorList, err.Error())
			continue
		}
		fmt.Fprintf(tw, fmt.Sprintf("F%03d %s\t: %%s\n", i, desc), str)
	}

	tw.Flush()

	if len(errorList) > 0 {
		fmt.Fprintf(w, "\nUnpacking Errors:\n")
		for _, err := range errorList {
			fmt.Fprintf(w, "- %s:\n", err)
		}

		return fmt.Errorf("displaying fields: %s", strings.Join(errorList, ","))
	}

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
