package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

type InRule struct {
	values []interface{}
}

func (m *InRule) Init(kind reflect.Kind, args []string) error {
	m.values = make([]interface{}, len(args))
	if !m.supports(kind) {
		return ErrSupportArgType
	}
	if len(args) == 0 {
		return ErrWrongArgsList
	}
	for i, arg := range args {
		var (
			v   interface{}
			err error
		)
		if kind >= reflect.Int && kind <= reflect.Int64 {
			v, err = strconv.ParseInt(arg, 10, 64)
		}
		if kind >= reflect.Uint && kind <= reflect.Uint64 {
			v, err = strconv.ParseInt(arg, 10, 64)
		}
		if kind >= reflect.Float32 && kind <= reflect.Float64 {
			v, err = strconv.ParseFloat(args[0], 64)
		}
		if kind == reflect.String {
			v = arg
		}
		if err != nil {
			return err
		}
		m.values[i] = v
	}
	return nil
}

func (m InRule) Check(val reflect.Value) error {
	kind := val.Kind()
	if !m.supports(kind) {
		return ErrSupportArgType
	}
	for _, iv := range m.values {
		if kind >= reflect.Int && kind <= reflect.Uint64 {
			if iv == val.Int() {
				return nil
			}
		}
		if kind >= reflect.Float32 && kind <= reflect.Float64 {
			if iv == val.Float() {
				return nil
			}
		}
		if kind == reflect.String {
			if iv == val.String() {
				return nil
			}
		}
	}
	return Invalid{Code: "in", Err: fmt.Errorf("`%v` is not in required set %v", val.Interface(), m.values)}
}

func (m InRule) supports(k reflect.Kind) bool {
	return k == reflect.String || (k >= reflect.Int && k <= reflect.Float64 && k != reflect.Uintptr)
}

func NewInRule() Rule {
	return &InRule{}
}
