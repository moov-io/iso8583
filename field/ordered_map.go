package field

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/moov-io/iso8583/utils"
)

// Custom type to sort keys in resulting JSON
type OrderedMap map[string]Field

func (om OrderedMap) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(om))
	for k := range om {
		keys = append(keys, k)
	}
	sort.Sort(sortImpl(keys))

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

type sortImpl []string

func (a sortImpl) Len() int      { return len(a) }
func (a sortImpl) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortImpl) Less(i, j int) bool {
	numLeft, err := strconv.ParseUint(a[i], 10, 0)
	if err != nil {
		return a[i] < a[j]
	}
	numRight, err := strconv.ParseUint(a[j], 10, 0)
	if err != nil {
		return a[i] < a[j]
	}
	return numLeft < numRight
}
