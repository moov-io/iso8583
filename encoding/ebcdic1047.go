package encoding

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
	xencoding "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// EBCDIC1047 is an encoder for EBCDIC characters using IBM Code Page 1047.
var EBCDIC1047 Encoder = &ebcdic1047Encoder{
	encoder: charmap.CodePage1047.NewEncoder(),
	decoder: charmap.CodePage1047.NewDecoder(),
}

type ebcdic1047Encoder struct {
	encoder *xencoding.Encoder
	decoder *xencoding.Decoder
}

func (e ebcdic1047Encoder) Encode(data []byte) ([]byte, error) {
	bytes, err := e.encoder.Bytes(data)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to encode EBCDIC")
	}
	return bytes, nil
}

func (e ebcdic1047Encoder) Decode(data []byte, length int) ([]byte, int, error) {
	if len(data) < length {
		return nil, 0, fmt.Errorf(
			"not enough data to decode. expected len %d, got %d", length, len(data),
		)
	}

	data = data[:length]
	out, err := e.decoder.Bytes(data)
	if err != nil {
		return nil, 0, utils.NewSafeError(err, "failed to decode EBCDIC")
	}
	return out, length, nil
}
