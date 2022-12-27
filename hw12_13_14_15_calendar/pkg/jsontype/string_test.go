package jsontype

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Customtype_String_Unmarshal_And_Marshal(t *testing.T) {
	type InputString struct {
		String String `json:"string"`
		Int    String `json:"int"`
		Float  String `json:"float"`
		Bool   String `json:"bool"`
	}
	var err error
	var val InputString
	var jsonString string
	var input []InputString

	jsonString = `[
		{"string" : "text", "int" : 10, "float" : 10.5, "bool" : false},
		{"string" : "text", "int" : "10", "float" : "10.5", "bool" : "false"}
	]`

	err = json.Unmarshal([]byte(jsonString), &input)
	require.NoError(t, err)
	for _, val = range input {
		require.Equal(t, val.String, String("text"))
		require.Equal(t, val.Int, String("10"))
		require.Equal(t, val.Float, String("10.5"))
		require.Equal(t, val.Bool, String("false"))
	}
	expectedJSON := `[{
		"string":"text",
		"int":"10",
		"float":"10.5",
		"bool":"false"
	},{
		"string":"text",
		"int":"10",
		"float":"10.5",
		"bool":"false"
	}]`
	makeJSON, err := json.Marshal(input)
	require.NoError(t, err)
	require.JSONEq(t, string(makeJSON), expectedJSON)
}
