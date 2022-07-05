package specs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/moov-io/iso8583"
)

// CreateFromJsonFile returns a MessageSpec generated from the input file
func CreateFromJsonFile(path string) (*iso8583.MessageSpec, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", path, err)
	}
	defer fd.Close()

	raw, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	return Builder.ImportJSON(raw)
}
