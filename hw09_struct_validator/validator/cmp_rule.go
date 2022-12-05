package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

// CmpRule правило, которое проверяет числовое значение функцией CmpFn
// считаем что самый общий числовой тип float64
type CmpRule struct {
	CmpFn     func(float64, float64) bool
	ErrFormat string
	value     float64 // самый общий числовой тип
}

func (m *CmpRule) Init(kind reflect.Kind, args []string) error {
	if !m.supports(kind) {
		return ErrSupportArgType
	}
	if len(args) != 1 {
		return ErrWrongArgsList
	}
	v, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return err
	}
	m.value = v

	return nil
}
func (m CmpRule) Check(val reflect.Value) error {
	var v float64
	t := val.Kind()
	if t >= reflect.Int && t <= reflect.Uint64 {
		v = float64(val.Int())
		if !m.CmpFn(v, m.value) {
			return fmt.Errorf(m.ErrFormat, m.value, v)
		}
		return nil
	}
	if t >= reflect.Float32 && t <= reflect.Float64 {
		v = val.Float()
		if !m.CmpFn(v, m.value) {
			return fmt.Errorf(m.ErrFormat, m.value, v)
		}
		return nil
	}
	return ErrSupportArgType
}
func (m CmpRule) supports(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Float64 && k != reflect.Uintptr
}

func NewMinRule() Rule {
	return &CmpRule{
		CmpFn: func(v float64, m float64) bool {
			return v >= m
		},
		ErrFormat: "expected value must be not less then %f, got %f",
	}
}
func NewMaxRule() Rule {
	return &CmpRule{
		CmpFn: func(v float64, m float64) bool {
			return v <= m
		},
		ErrFormat: "expected value must be not greater then %f, got %f",
	}
}
