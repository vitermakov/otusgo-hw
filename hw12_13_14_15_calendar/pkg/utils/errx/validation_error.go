package errx

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("{%s} - %s", ve.Field, ve.Err.Error())
}

type ValidationErrors []ValidationError

func (vs ValidationErrors) Error() string {
	if len(vs) == 0 {
		return ""
	}
	err := strings.Builder{}
	for i, ve := range vs {
		err.WriteString(ve.Error())
		if i > 0 {
			err.WriteString("; ")
		}
	}
	return err.String()
}

func (vs *ValidationErrors) Add(ve ValidationError) {
	*vs = append(*vs, ve)
}

func (vs ValidationErrors) Empty() bool {
	return len(vs) == 0
}
