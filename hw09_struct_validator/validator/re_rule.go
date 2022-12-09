package validator

import (
	"fmt"
	"reflect"
	"regexp"
)

// ReRule проверка строки по регулярному выражению.
type ReRule struct {
	re regexp.Regexp // скомпилированное рег.выражение.
}

func (m *ReRule) Init(kind reflect.Kind, args []string) error {
	if !m.supports(kind) {
		return ErrSupportArgType
	}
	if len(args) != 1 {
		return ErrWrongArgsList
	}
	r, err := regexp.Compile(args[0])
	if err != nil {
		return err
	}
	m.re = *r

	return nil
}

func (m ReRule) Check(val reflect.Value) error {
	if !m.supports(val.Kind()) {
		return ErrSupportArgType
	}
	if !m.re.Match([]byte(val.String())) {
		return Invalid{
			Code: "re",
			Err:  fmt.Errorf("value `%s` not matching spcified pattern `%s`", val.String(), m.re.String()),
		}
	}
	return nil
}

func (m ReRule) supports(k reflect.Kind) bool {
	return k == reflect.String
}

func NewReRule() Rule {
	return &ReRule{}
}
