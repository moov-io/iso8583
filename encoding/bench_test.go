package encoding

import "testing"

var (
	benchmarkASCII  = []byte("1234567890123456")
	benchmarkEBCDIC = []byte{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6}
	benchmarkLBCD   = []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56}
	benchmarkHex    = []byte{0xAB, 0xCD, 0xEF, 0x01, 0x23, 0x45, 0x67, 0x89}
	benchmarkTag    = []byte{0x9F, 0x02}
)

func BenchmarkASCIIEncode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := ASCII.Encode(benchmarkASCII); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkASCIIDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := ASCII.Decode(benchmarkASCII, len(benchmarkASCII)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBinaryDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := Binary.Decode(benchmarkASCII, len(benchmarkASCII)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEBCDICEncode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := EBCDIC.Encode(benchmarkASCII); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEBCDICDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := EBCDIC.Decode(benchmarkEBCDIC, len(benchmarkEBCDIC)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLBCDEncode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := LBCD.Encode(benchmarkASCII); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLBCDDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := LBCD.Decode(benchmarkLBCD, len(benchmarkLBCD)*2); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBytesToASCIIHexEncode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := BytesToASCIIHex.Encode(benchmarkHex); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkASCIIHexToBytesDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := ASCIIHexToBytes.Decode(benchmarkHex, len(benchmarkHex)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBerTLVTagDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := BerTLVTag.Decode(benchmarkTag, 0); err != nil {
			b.Fatal(err)
		}
	}
}
