package jsontype

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Customtype_Float_Unmarshal_And_Marshal(t *testing.T) {

	type InputFloat struct {
		Float32 Float32 `json:"float32"`
		Float64 Float64 `json:"float64"`
	}

	var err error
	var val InputFloat
	var jsonString string
	var input []InputFloat

	jsonString = `[
		{"float32": 1,  "float64": 2 },
		{"float32":"1", "float64":"2"}
	]`

	if err = json.Unmarshal([]byte(jsonString), &input); assert.NoError(t, err) {
		for _, val = range input {
			assert.Equal(t, val.Float32, Float32(1))
			assert.Equal(t, val.Float64, Float64(2))
		}
		if make_json, err := json.Marshal(input); assert.NoError(t, err) {
			assert.Equal(t, string(make_json), `[{"float32":1,"float64":2},{"float32":1,"float64":2}]`)
		}
	}

	jsonString = `[
		{"float32": 1.5 , "float64": 2.5 },
		{"float32":"1.5", "float64":"2.5"}
	]`

	if err = json.Unmarshal([]byte(jsonString), &input); assert.NoError(t, err) {
		for _, val = range input {
			assert.Equal(t, val.Float32, Float32(1.5))
			assert.Equal(t, val.Float64, Float64(2.5))
		}
		if make_json, err := json.Marshal(input); assert.NoError(t, err) {
			assert.Equal(t, string(make_json), `[{"float32":1.5,"float64":2.5},{"float32":1.5,"float64":2.5}]`)
		}
	}
}

func Test_Customtype_Float_Unmarshal_Error(t *testing.T) {

	type InputFloat struct {
		Float32 Float32 `json:"float32"`
		Float64 Float64 `json:"float64"`
	}

	var err error
	var jsonString string
	var input []InputFloat

	jsonString = `[
		{"float32": 1,  "float64": "a1" }
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"float32": "a1",  "float64": 1 }
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"float32": false,  "float64": 1 }
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value false invalid`)
	}

	jsonString = `[
		{"float32": 1,  "float64": true }
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value true invalid`)
	}
}
