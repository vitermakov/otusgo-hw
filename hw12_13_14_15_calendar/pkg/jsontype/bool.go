package jsontype

import (
	"errors"
	"strings"
)

type Bool bool

func (value *Bool) UnmarshalJSON(data []byte) (err error) {
	switch strings.ToLower(string(data)) {
	case `"true"`, `true`, `"1"`, `1`, `"yes"`, `"y"`, `"on"`:
		*value = Bool(true)
		break
	case `"false"`, `false`, `"0"`, `0`, `""`, `"no"`, `"n"`, `"off"`:
		*value = Bool(false)
		break
	default:
		err = errors.New("Value " + string(data) + " invalid")
		break
	}
	return
}