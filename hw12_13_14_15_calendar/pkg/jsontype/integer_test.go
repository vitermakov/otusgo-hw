package jsontype

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Customtype_Integer_Unmarshal_And_Marshal(t *testing.T) {

	type InputInteger struct {
		Int    Int    `json:"int"`
		Int32  Int32  `json:"int32"`
		Int64  Int64  `json:"int64"`
		Uint   Uint   `json:"uint"`
		Uint32 Uint32 `json:"uint32"`
		Uint64 Uint64 `json:"uint64"`
	}

	var err error
	var val InputInteger
	var jsonString string

	var input []InputInteger

	jsonString = `[
		{"int":  -1 , "int32":  32 , "int64":  64 , "uint":  1 , "uint32":  32 , "uint64":  64 },
		{"int": "-1", "int32": "32", "int64": "64", "uint": "1", "uint32": "32", "uint64": "64"}
	]`

	if err = json.Unmarshal([]byte(jsonString), &input); assert.NoError(t, err) {
		for _, val = range input {
			assert.Equal(t, val.Int, Int(-1))
			assert.Equal(t, val.Int32, Int32(32))
			assert.Equal(t, val.Int64, Int64(64))
			assert.Equal(t, val.Uint, Uint(1))
			assert.Equal(t, val.Uint32, Uint32(32))
			assert.Equal(t, val.Uint64, Uint64(64))
		}

		if make_json, err := json.Marshal(input); assert.NoError(t, err) {
			assert.Equal(t, string(make_json), `[{"int":-1,"int32":32,"int64":64,"uint":1,"uint32":32,"uint64":64},{"int":-1,"int32":32,"int64":64,"uint":1,"uint32":32,"uint64":64}]`)
		}
	}
}

func Test_Customtype_Integer_Unmarshal_Error(t *testing.T) {

	type InputInteger struct {
		Int    Int    `json:"int"`
		Int32  Int32  `json:"int32"`
		Int64  Int64  `json:"int64"`
		Uint   Uint   `json:"uint"`
		Uint32 Uint32 `json:"uint32"`
		Uint64 Uint64 `json:"uint64"`
	}

	var err error
	var jsonString string
	var input []InputInteger

	jsonString = `[
		{"int": "a1", "int32": 1, "int64": 1, "uint": 1, "uint32": 1, "uint64": 1}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"int": 1, "int32": "a1", "int64": 1, "uint": 1, "uint32": 1, "uint64": 1}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"int": 1, "int32": 1, "int64": "a1", "uint": 1, "uint32": 1, "uint64": 1}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"int": 1, "int32": 1, "int64": 1, "uint": "a1", "uint32": 1, "uint64": 1}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"int": 1, "int32": 1, "int64": 1, "uint": 1, "uint32": "a1", "uint64": 1}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"int": 1, "int32": 1, "int64": 1, "uint": 1, "uint32": 1, "uint64": "a1"}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value "a1" invalid`)
	}

	jsonString = `[
		{"int": 1, "int32": 1, "int64": 1, "uint": 1, "uint32": 1, "uint64": true}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value true invalid`)
	}

	jsonString = `[
		{"int": 1, "int32": 1, "int64": 1, "uint": 1, "uint32": 1, "uint64": false}
	]`
	if err = json.Unmarshal([]byte(jsonString), &input); assert.Error(t, err) {
		assert.Equal(t, err.Error(), `Value false invalid`)
	}
}
