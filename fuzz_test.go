package iso8583_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
)

func FuzzUnpackPackSpec87(f *testing.F) {
	populateISOCorpus(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64*1024 {
			t.Skip()
		}

		msg := iso8583.NewMessage(iso8583.Spec87)
		if err := msg.Unpack(data); err != nil {
			return
		}
		packed, err := msg.Pack()
		if err != nil {
			return
		}
		// Re-unpack packed bytes — must not panic
		msg2 := iso8583.NewMessage(iso8583.Spec87)
		_ = msg2.Unpack(packed)
		_, _ = msg.MarshalJSON()
	})
}

func FuzzUnpackPackSpec87ASCII(f *testing.F) {
	populateISOCorpus(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64*1024 {
			t.Skip()
		}

		msg := iso8583.NewMessage(specs.Spec87ASCII)
		if err := msg.Unpack(data); err != nil {
			return
		}
		_, _ = msg.Pack()
		_, _ = msg.MarshalJSON()
	})
}

func FuzzMessageScanner(f *testing.F) {
	populateISOCorpus(f)

	// Common field IDs to probe with the forward-only scanner
	fieldIDs := []int{0, 2, 3, 4, 7, 11, 39, 41, 49}

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64*1024 {
			t.Skip()
		}

		scanner := iso8583.NewMessageScanner(iso8583.Spec87, data)
		for _, id := range fieldIDs {
			_, _ = scanner.ScanField(id)
		}
	})
}

func populateISOCorpus(f *testing.F) {
	f.Helper()

	f.Add([]byte{})
	f.Add([]byte("0800"))
	f.Add([]byte("080082200000000000000400000000000000123456789012"))

	roots := []string{
		filepath.Join("test", "testdata"),
		filepath.Join("test", "fuzz-reader", "corpus"),
		filepath.Join("testdata"),
	}
	for _, root := range roots {
		_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			if ext == ".dat" || ext == "" || ext == ".bin" {
				bs, err := os.ReadFile(path)
				if err != nil || len(bs) > 64*1024 {
					return nil
				}
				f.Add(bs)
			}
			return nil
		})
	}
}
