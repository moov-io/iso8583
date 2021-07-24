package specs

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {

	asciiJson, err := Builder.ExportJSON(Spec87ASCII)
	require.NoError(t, err)

	asciiSpec, err := Builder.ImportJSON(asciiJson)
	require.NoError(t, err)

	require.Equal(t, true, reflect.DeepEqual(Spec87ASCII, asciiSpec))

	hexJson, err := Builder.ExportJSON(Spec87Hex)
	require.NoError(t, err)

	hexSpec, err := Builder.ImportJSON(hexJson)
	require.NoError(t, err)

	require.Equal(t, true, reflect.DeepEqual(Spec87Hex.Name, hexSpec.Name))

}

func TestExampleJSONSpec(t *testing.T) {
	asciiJson, err := ioutil.ReadFile("../examples/specs/spec87ascii.json")
	require.NoError(t, err)

	asciiSpec, err := Builder.ImportJSON(asciiJson)
	require.NoError(t, err)

	require.Equal(t, true, reflect.DeepEqual(Spec87ASCII, asciiSpec))
}
