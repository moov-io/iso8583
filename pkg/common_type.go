package pkg

// general element type for all of the data representation attributes
type CommonType struct {
	Type     string
	Length   int
	Fixed    bool
	Format   string
	Encoding string
}

func (t *CommonType) Validate() error {
	return nil
}
