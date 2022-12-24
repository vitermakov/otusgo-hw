package jsontype

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Customtype_Bool_Unmarshal_And_Marshal(t *testing.T) {

	type InputBool struct {
		BoolTrue  Bool `json:"booltrue"`
		BoolFalse Bool `json:"boolfalse"`
	}

	var err error
	var val InputBool
	var jsonString string
	var input []InputBool

	jsonString = `[
		{"booltrue": true   , "boolfalse": false},
		{"booltrue": "true" , "boolfalse": "false"},
		{"booltrue":  1     , "boolfalse":  0},
		{"booltrue": "1"    , "boolfalse": "0"},
		{"booltrue": "on"   , "boolfalse": "off"},
		{"booltrue": "Yes"  , "boolfalse": "No"},
		{"booltrue": "yes"  , "boolfalse": "no"},
		{"booltrue": "y"    , "boolfalse": "n"}
	]`

	if err = json.Unmarshal([]byte(jsonString), &input); assert.NoError(t, err) {
		for _, val = range input {
			assert.Equal(t, val.BoolTrue, Bool(true))
			assert.Equal(t, val.BoolFalse, Bool(false))
		}

		if make_json, err := json.Marshal(input); assert.NoError(t, err) {
			assert.Equal(t, string(make_json), `[{"booltrue":true,"boolfalse":false},{"booltrue":true,"boolfalse":false},{"booltrue":true,"boolfalse":false},{"booltrue":true,"boolfalse":false},{"booltrue":true,"boolfalse":false},{"booltrue":true,"boolfalse":false},{"booltrue":true,"boolfalse":false},{"booltrue":true,"boolfalse":false}]`)
		}
	}
}

func Test_Customtype_Bool_Unmarshal_Error(t *testing.T) {

	type InputBool struct {
		BoolTrue  Bool `json:"booltrue"`
		BoolFalse Bool `json:"boolfalse"`
	}

	var err error
	var jsonString string
	var input []InputBool

	jsonString = `[
		{"booltrue": "true_"   , "boolfalse": false},
		{"booltrue": true   , "boolfalse": "false_"}
	]`

	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "true_" invalid`)
	}

	jsonString = `[
		{"booltrue": true   , "boolfalse": false},
		{"booltrue": true   , "boolfalse": "false_"}
	]`

	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "false_" invalid`)
	}
}
