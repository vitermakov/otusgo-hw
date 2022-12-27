package jsontype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Customtype_Base_UnmarshalJSON(t *testing.T) {
	var err error
	var jsonString []string
	var input int64

	jsonString = []string{
		`1`,
		`"1"`,
	}

	for _, val := range jsonString {
		if err = UnmarshalJSON([]byte(val), &input); assert.NoError(t, err) {
			assert.Equal(t, input, int64(1))
		}
	}
}
