package errx

import (
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return ve.Err.Error()
}

type ValidationErrors []ValidationError

func (vs ValidationErrors) Error() string {
	if len(vs) == 0 {
		return ""
	}
	err := strings.Builder{}
	for _, i := range vs {
		str := i.Error()
		err.WriteString(str + "\n")
	}
	return err.String()
}

func (vs *ValidationErrors) Add(ve ValidationError) {
	*vs = append(*vs, ve)
}

func (vs ValidationErrors) Empty() bool {
	return len(vs) == 0
}
