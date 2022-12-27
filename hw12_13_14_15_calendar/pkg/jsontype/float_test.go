package jsontype

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
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

	err = json.Unmarshal([]byte(jsonString), &input)
	require.NoError(t, err)
	for _, val = range input {
		require.Equal(t, val.Float32, Float32(1))
		require.Equal(t, val.Float64, Float64(2))
	}
	makeJSON, err := json.Marshal(input)
	require.NoError(t, err)
	require.JSONEq(t, string(makeJSON), `[{"float32":1,"float64":2},{"float32":1,"float64":2}]`)

	jsonString = `[
		{"float32": 1.5 , "float64": 2.5 },
		{"float32":"1.5", "float64":"2.5"}
	]`

	err = json.Unmarshal([]byte(jsonString), &input)
	require.NoError(t, err)
	for _, val = range input {
		require.Equal(t, val.Float32, Float32(1.5))
		require.Equal(t, val.Float64, Float64(2.5))
	}
	makeJSON, err = json.Marshal(input)
	require.NoError(t, err)
	require.JSONEq(t, string(makeJSON), `[{"float32":1.5,"float64":2.5},{"float32":1.5,"float64":2.5}]`)
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
	err = json.Unmarshal([]byte(jsonString), &input)
	require.Error(t, err)
	require.Equal(t, err.Error(), `Value "a1" invalid`)

	jsonString = `[
		{"float32": "a1",  "float64": 1 }
	]`
	err = json.Unmarshal([]byte(jsonString), &input)
	require.Error(t, err)
	require.Equal(t, err.Error(), `Value "a1" invalid`)

	jsonString = `[
		{"float32": false,  "float64": 1 }
	]`
	err = json.Unmarshal([]byte(jsonString), &input)
	require.Error(t, err)
	require.Equal(t, err.Error(), `Value false invalid`)

	jsonString = `[
		{"float32": 1,  "float64": true }
	]`
	err = json.Unmarshal([]byte(jsonString), &input)
	require.Error(t, err)
	require.Equal(t, err.Error(), `Value true invalid`)
}
