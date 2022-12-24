package jsontype

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
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

	if err = json.Unmarshal([]byte(jsonString), &input); assert.NoError(t, err) {
		for _, val = range input {
			assert.Equal(t, val.String, String("text"))
			assert.Equal(t, val.Int, String("10"))
			assert.Equal(t, val.Float, String("10.5"))
			assert.Equal(t, val.Bool, String("false"))
		}

		if make_json, err := json.Marshal(input); assert.NoError(t, err) {
			assert.Equal(t, string(make_json), `[{"string":"text","int":"10","float":"10.5","bool":"false"},{"string":"text","int":"10","float":"10.5","bool":"false"}]`)
		}
	}
}
