package network

import (
	"io"
)

/*
Header is the interface for network header

All messages between client/server have a message length header. In some cases
it can be a 4 bytes ASCII or 2 bytes BCD encoded length.

Here is how you can use it to read message from net.Conn:

	func handleRequest(conn net.Conn) {
		header := network.NewBCD2BytesHeader()
		_, err := header.ReadFrom(conn)
		if err != nil {
			fmt.Printf("Reading header: %w\n", err)
		}

		// Make a buffer to hold message
		buf := make([]byte, header.Length())
		// Read the incoming connection into the buffer.
		read, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		if reqLen != header.Length() {
			fmt.Println("Expected to read %d bytes, read %d bytes", header.Length(), read)
		}

		message := iso8583.NewMessage(iso8583.Spec87)
		message.Unpack(buf)

		// work with the message
	}

This is how you can write header into the net.Conn:

	header := network.NewBCD2BytesHeader()
	packed, err := message.Pack()
	if err != nil {
		// handle error
	}
	header.SetLength(len(packed))
	_, err = header.WriteTo(conn)
	if err != nil {
		// handle error
	}
	n, err := conn.Write(packed)
	if err != nil {
		// handle error
	}

*/

// Header is network header interface to write/read encoded message legnth
type Header interface {
	// WriteTo encoded length into Writer
	WriteTo(w io.Writer) (int, error)

	// ReadFrom reads header (N bytes) from the Reader
	ReadFrom(r io.Reader) (int, error)

	// SetLength sets the length of the message
	SetLength(length int)

	// Length returns the length of the message
	Length() int
}
