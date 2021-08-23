package bcd

import (
	// "bytes"
	// "fmt"
	"io"
)

// Writer encodes input data and writes it to underlying io.Writer.
// Please pay attention that due to ambiguity of encoding process
// (encoded octet may indicate the end of data by using the filler
// nibble) Writer will not write odd remainder of the encoded input
// data if any until the next octet is observed.
type Writer struct {
	*Encoder
	dst  io.Writer
	err  error
	word []byte
}

// NewWriter creates new Writer with underlying io.Writer.
func (enc *Encoder) NewWriter(wr io.Writer) *Writer {
	return &Writer{enc, wr, nil, make([]byte, 0, 2)}
}

// Write implements io.Writer interface.
func (w *Writer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	// if we have remaining byte from previous run
	// join it with one of new input and encode
	if len(w.word) == 1 {
		x := append(w.word, p[0])
		_, b, err := w.pack(x)
		if err != nil {
			return 0, err
		}
		if _, err = w.dst.Write([]byte{b}); err != nil {
			return 0, err
		}
		w.word = w.word[:0]
		n += 1
	}

	// encode even number of bytes
	for len(p[n:]) >= 2 {
		_, b, err := w.pack(p[n : n+2])
		if err != nil {
			return n, err
		}
		if _, err = w.dst.Write([]byte{b}); err != nil {
			return n, err
		}
		n += 2
	}

	// save remainder
	if len(p[n:]) > 0 { // == 1
		w.word = append(w.word, p[n])
		n += 1
	}

	return
}

// Encodes all backlogged data to underlying Writer.  If number of
// bytes is odd, the padding fillers will be applied. Because of this
// the main usage of Flush is right before stopping Write()-ing data
// to properly finalize the encoding process.
func (w *Writer) Flush() error {
	if len(w.word) == 0 {
		return nil
	}
	n, b, err := w.pack(w.word)
	w.word = w.word[:0]
	if err != nil {
		// panic("hell")
		return err
	}
	if n == 0 {
		return nil
	}
	_, err = w.dst.Write([]byte{b})
	return err
}

// Buffered returns the number of bytes stored in backlog awaiting for
// its pair.
func (w *Writer) Buffered() int {
	return len(w.word)
}
