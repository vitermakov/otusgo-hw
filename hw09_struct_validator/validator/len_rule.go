package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

// CmpRule правило, которое проверяет числовое значение функцией CmpFn
// считаем что самый общий числовой тип float64
type LenRule struct {
	CmpFn     func(float64, float64) bool
	ErrFormat string
	lenTest   int
}

func (m *LenRule) Init(kind reflect.Kind, args []string) error {
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
	if v < 0 {
		return fmt.Errorf("len can't be negative, got %d", v)
	}
	m.lenTest = int(v)

	return nil
}
func (m LenRule) Check(val reflect.Value) error {
	if val.Kind() == reflect.String {
		if val.Len() != m.lenTest {
			return fmt.Errorf("expected length %d, got %d", m.lenTest, val.Len())
		}
		return nil
	}
	return ErrSupportArgType
}
func (m LenRule) supports(k reflect.Kind) bool {
	return k == reflect.String
}
func NewLenRule() Rule {
	return &LenRule{}
}
