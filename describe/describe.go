package describe

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/moov-io/iso8583"
)

func Message(w io.Writer, message *iso8583.Message) error {
	fmt.Fprintf(w, "ISO 8583 Message\n****************\n")

	mti, err := message.GetMTI()
	if err != nil {
		return fmt.Errorf("getting MTI: %w", err)
	}
	printField("MTI", mti)

	bitmapRaw, err := message.Bitmap().Bytes()
	if err != nil {
		return fmt.Errorf("getting bitmap bytes: %w", err)
	}
	printField("Bitmap", strings.ToUpper(hex.EncodeToString(bitmapRaw)))

	bitmap, err := message.Bitmap().String()
	if err != nil {
		return fmt.Errorf("getting bitmap: %w", err)
	}
	printField("Bitmap bits", bitmap)

	return nil
}

var lableLength = 30

func printField(name string, value interface{}) {
	padding := strings.Repeat(".", lableLength-len(name))
	fmt.Printf("%s%s: %v\n", name, padding, value)
}
