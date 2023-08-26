package field

import (
	"encoding/json"
	"testing"
)

var (
	testOrderedMap = OrderedMap{
		"0":   &String{value: "0100"},
		"1":   &String{value: "D00080000000000000000000002000000000000000000000"},
		"107": &String{value: "102"},
		"17":  &String{value: "101"},
		"2":   &String{value: "4242424242424242"},
		"4":   &String{value: "100"},
		"a":   &String{value: "555"},
		"z":   &String{value: "555"},
	}
)

func TestOrderedMap_MarshalJSON(t *testing.T) {
	expected := `{"0":"0100","1":"D00080000000000000000000002000000000000000000000","2":"4242424242424242","4":"100","17":"101","107":"102","a":"555","z":"555"}`
	data, err := json.Marshal(testOrderedMap)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != expected {
		t.Errorf("Marshalled value should be \n\t%s\ninstead of \n\t%s", expected, string(data))
	}
}

func BenchmarkOrderedMap_MarshalJSON(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(testOrderedMap)
		if err != nil {
			b.Fatal(err)
		}
	}
}
