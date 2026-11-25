// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package examples_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/network"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

// ascii3Header encodes the message body length as three ASCII digits.
type ascii3Header struct {
	length int
}

var _ network.Header = (*ascii3Header)(nil)

func (h *ascii3Header) SetLength(length int) {
	h.length = length
}

func (h *ascii3Header) Length() int {
	return h.length
}

func (h *ascii3Header) WriteTo(w io.Writer) (int, error) {
	if h.length < 0 || h.length > 999 {
		return 0, fmt.Errorf("message length %d is outside the range 0 to 999", h.length)
	}
	return fmt.Fprintf(w, "%03d", h.length)
}

func (h *ascii3Header) ReadFrom(r io.Reader) (int, error) {
	var digits [3]byte
	n, err := io.ReadFull(r, digits[:])
	if err != nil {
		return n, err
	}

	length := 0
	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			return n, fmt.Errorf("invalid ASCII digit %q in length header", digit)
		}
		length = length*10 + int(digit-'0')
	}
	h.length = length
	return n, nil
}

func Example_customNetworkHeader() {
	msg := iso8583.NewMessage(specs.Spec87ASCII)
	msg.MTI("0800")

	body, err := msg.Pack()
	if err != nil {
		panic(err)
	}

	var frame bytes.Buffer
	var outgoing ascii3Header
	outgoing.SetLength(len(body))
	if _, err := outgoing.WriteTo(&frame); err != nil {
		panic(err)
	}
	frame.Write(body)

	var incoming ascii3Header
	if _, err := incoming.ReadFrom(&frame); err != nil {
		panic(err)
	}
	received := make([]byte, incoming.Length())
	if _, err := io.ReadFull(&frame, received); err != nil {
		panic(err)
	}

	fmt.Println("Body matches:", bytes.Equal(body, received))
	// Output:
	// Body matches: true
}

func TestASCII3HeaderWriteTo(t *testing.T) {
	for _, tc := range []struct {
		name   string
		length int
		want   string
	}{
		{"zero", 0, "000"},
		{"leading zero", 27, "027"},
		{"maximum", 999, "999"},
		{"negative", -1, ""},
		{"overflow", 1000, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var header ascii3Header
			header.SetLength(tc.length)
			require.Equal(t, tc.length, header.Length())

			var buf bytes.Buffer
			n, err := header.WriteTo(&buf)
			if tc.want == "" {
				require.ErrorContains(t, err, "outside the range 0 to 999")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, len(tc.want), n)
			require.Equal(t, tc.want, buf.String())
		})
	}
}

type partialHeaderWriter struct{ err error }

func (w partialHeaderWriter) Write([]byte) (int, error) {
	return 1, w.err
}

func TestASCII3HeaderWriteToPartialWrite(t *testing.T) {
	cause := errors.New("write interrupted")
	var header ascii3Header
	header.SetLength(27)
	n, err := header.WriteTo(partialHeaderWriter{err: cause})
	require.ErrorIs(t, err, cause)
	require.Equal(t, 1, n)
}

func TestASCII3HeaderReadFrom(t *testing.T) {
	t.Run("leaves the message body unread", func(t *testing.T) {
		stream := strings.NewReader("003ABC")
		var header ascii3Header
		n, err := header.ReadFrom(stream)
		require.NoError(t, err)
		require.Equal(t, 3, n)
		require.Equal(t, 3, header.Length())
		body, err := io.ReadAll(stream)
		require.NoError(t, err)
		require.Equal(t, []byte("ABC"), body)
	})

	t.Run("truncated header", func(t *testing.T) {
		var header ascii3Header
		header.SetLength(27)
		n, err := header.ReadFrom(strings.NewReader("02"))
		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
		require.Equal(t, 2, n)
		require.Equal(t, 27, header.Length())
	})

	t.Run("invalid digit", func(t *testing.T) {
		var header ascii3Header
		header.SetLength(27)
		n, err := header.ReadFrom(strings.NewReader("0A7"))
		require.ErrorContains(t, err, "invalid ASCII digit 'A'")
		require.Equal(t, 3, n)
		require.Equal(t, 27, header.Length())
	})
}
