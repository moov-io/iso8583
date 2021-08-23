package bcd

import (
	"bytes"
	// "fmt"
	"io"
)

// Reader reads encoded BCD data from underlying io.Reader and decodes
// them. Please pay attention that due to ambiguity of encoding
// process (encoded octet may indicate the end of data by using the
// filler nibble) the last input octet is not decoded until the next
// input octet is observed or until underlying io.Reader returns
// error.
type Reader struct {
	*Decoder
	src io.Reader
	err error
	buf bytes.Buffer
	out []byte
}

// NewReader creates new Reader with underlying io.Reader.
func (dec *Decoder) NewReader(rd io.Reader) *Reader {
	return &Reader{dec, rd, nil, bytes.Buffer{}, []byte{}}
}

// Read implements io.Reader interface.
func (r *Reader) Read(p []byte) (n int, err error) {
	buf := &r.buf

	// return previously decoded data first
	backlog := copy(p[n:], r.out)
	r.out = r.out[backlog:]
	n += backlog
	if len(p) == n {
		return
	}

	if x := EncodedLen(len(p)); r.err == nil {
		// refill on data
		_, r.err = io.CopyN(buf, r.src, int64(x+1))
	}

	if r.err != nil && buf.Len() == 0 {
		// underlying Reader gives no data,
		// buffer is also empty, we're done
		return n, r.err
	}

	// decoding buffer
	w := make([]byte, 2)

	// no error yet, we have some data to decode;
	// decoding until the only byte is left in buffer
	for buf.Len() > 1 && n < len(p) {
		b, _ := buf.ReadByte()
		wid, end, err := r.unpack(w, b)
		if err != nil {
			return n, err
		}

		if end && !r.IgnoreFiller {
			err = ErrBadBCD
		}

		// fmt.Printf("copying '%c' '%c' - %d bytes\n", w[0], w[1], wid)
		cp := copy(p[n:], w[:wid])
		r.out = append(r.out, w[cp:wid]...)
		n += cp

		if err != nil {
			return n, err
		}
	}

	// last breath
	if buf.Len() == 1 && r.err != nil {
		b, _ := buf.ReadByte()
		wid, _, err := r.unpack(w, b)
		if err != nil {
			return n, err
		}

		// fmt.Printf("copying '%c' '%c' - %d bytes\n", w[0], w[1], wid)
		cp := copy(p[n:], w[:wid])
		r.out = append(r.out, w[cp:wid]...)
		n += cp
	}

	return
}
