package field

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/moov-io/iso8583/utils"
)

// Custom type to sort keys in resulting JSON
type OrderedMap map[string]Field

func (om OrderedMap) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(om))
	for k := range om {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for _, i := range keys {
		b, err := json.Marshal(om[i])
		if err != nil {
			return nil, utils.NewSafeError(err, "failed to JSON marshal field to bytes")
		}
		buf.WriteString(fmt.Sprintf("\"%v\":", i))
		buf.Write(b)

		// don't add "," if it's the last item
		if i == keys[len(keys)-1] {
			break
		}

		buf.Write([]byte{','})
	}
	buf.Write([]byte{'}'})

	return buf.Bytes(), nil
}
