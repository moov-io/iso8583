package prefixer

// type Prefixer interface {
// 	EncodeLength(maxLength, length int) ([]byte, error)
// 	DecodeLength(maxLength int, data []byte) (int, error)
// 	DecodedLength() int
// 	Inspect() string
// }

// var Fixed Prefixer = &nonePrefixer{}

// var None Prefixer = &nonePrefixer{}

// type nonePrefixer struct{}

// func (*nonePrefixer) EncodeLength(int) ([]byte, error) {
// 	return nil, nil
// }
// func (*nonePrefixer) DecodeLength([]byte) (int, error) {
// 	return -1, nil
// }
// func (*nonePrefixer) DecodedLength() int {
// 	return -1
// }
// func (*nonePrefixer) Inspect() string {
// 	return "None prefixer"
// }

type Prefixer interface {
	EncodeLength(length int) ([]byte, error)
	DecodeLength(data []byte) (int, error)
	DecodedLength() int
	Inspect() string
}

type varPrefixers struct {
	L    PrefixerBuilder
	LL   PrefixerBuilder
	LLL  PrefixerBuilder
	LLLL PrefixerBuilder
}

type PrefixerBuilder func(int) Prefixer
