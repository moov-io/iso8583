package spec

type MessageSpec struct {
	Fields map[int]Packer
}

type Packer interface {
	// Pack packs the data taking into account data encoding and data length.
	Pack(data []byte) ([]byte, error)

	// Unpack unpacks data taking into account data encoding and data length
	// it returns unpacked data and the number of bytes read
	Unpack(data []byte) ([]byte, int, error)
}
