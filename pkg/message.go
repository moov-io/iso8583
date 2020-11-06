package pkg

// message is structure for ISO 8583 message encode and decode
type Message struct {
	Mti            *Element
	Bitmap         *Element
	Elements       *DataElements
	Specifications *Specification
	SecondBitmap   bool
	ThirdBitmap    bool
}
