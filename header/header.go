package header

import (
	"io"
)

// Header is the interface for ISO Headers
//
// If iso8583 Message has a header then `message.Read` will read the header
// first, and after decoding the length from the header will read
// `header.Length()` bytes of the message and unpack read data. Similar happens
// for `message.Write` if header is set. It packs the message first, then
// prepends header with the length of the message to the packed message.
type Header interface {
	// Pack packs the length of the header
	Pack() ([]byte, error)

	// Read reads N bytes of the header from the Reader
	Read(reader io.Reader) (int, error)

	// SetLength sets the length of the message
	SetLength(length int)

	// Length returns the length of the message
	Length() int
}
