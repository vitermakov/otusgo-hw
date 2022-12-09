package validator

import (
	"fmt"
	"reflect"
	"regexp"
)

// reRule проверка строки по регулярному выражению.
type reRule struct {
	re regexp.Regexp // скомпилированное рег.выражение.
}

func (m *reRule) init(kind reflect.Kind, args []string) error {
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

func (m reRule) Check(val reflect.Value) error {
	if !m.supports(val.Kind()) {
		return ErrSupportArgType
	}
	if !m.re.Match([]byte(val.String())) {
		return Invalid{
			Code: "regexp",
			Err:  fmt.Errorf("value `%s` not matching spcified pattern `%s`", val.String(), m.re.String()),
		}
	}
	return nil
}

func (m reRule) supports(k reflect.Kind) bool {
	return k == reflect.String
}

func createReRule(kind reflect.Kind, args []string) (rule, error) {
	rule := &reRule{}
	if err := rule.init(kind, args); err != nil {
		return nil, err
	}
	return rule, nil
}
