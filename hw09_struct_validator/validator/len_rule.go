package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

// LenRule проверяет строковое поле на длину.
type lenRule struct {
	ErrFormat string
	lenTest   int
}

func (m *lenRule) init(kind reflect.Kind, args []string) error {
	if !m.supports(kind) {
		return ErrSupportArgType
	}
	if len(args) != 1 {
		return ErrWrongArgsList
	}
	v, err := strconv.ParseInt(args[0], 0, 64)
	if err != nil {
		return err
	}
	if v <= 0 {
		return fmt.Errorf("len must be positive, got %d", v)
	}
	m.lenTest = int(v)

	return nil
}

func (m lenRule) Check(val reflect.Value) error {
	if !m.supports(val.Kind()) {
		return ErrSupportArgType
	}
	if val.Len() != m.lenTest {
		return Invalid{Code: "cmp", Err: fmt.Errorf("expected length %d, got %d", m.lenTest, val.Len())}
	}
	return nil
}

func (m lenRule) supports(k reflect.Kind) bool {
	return k == reflect.String
}

func createLenRule(kind reflect.Kind, args []string) (rule, error) {
	rule := &lenRule{}
	if err := rule.init(kind, args); err != nil {
		return nil, err
	}
	return rule, nil
}
