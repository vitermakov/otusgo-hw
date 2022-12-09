package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrInputStructIsNull = errors.New("variable value is null")
	ErrInputNotStruct    = errors.New("variable value is not struct")
	ErrWrongArgsList     = errors.New("wrong argument list")
	ErrSupportArgType    = errors.New("unsupported arg type")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	fmt.Println(len(v))
	s := strings.Builder{}
	for _, err := range v {
		s.WriteString(err.Field + ": " + err.Err.Error() + "\n")
	}
	return s.String()
}

type Rules []Rule

type StructRules map[int]interface{}

func ValidateStruct(v interface{}) error {
	if v == nil {
		return ErrInputStructIsNull
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrInputNotStruct
	}
	typ := val.Type()
	key := ""
	if typ.Name() != "" {
		key = typ.Name() + ":" + typ.PkgPath()
	}
	rules, err := retrieveRules(val, key)
	if err != nil {
		return err
	}
	errs := checkStruct(val, rules, typ.Name())
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func retrieveRules(rStruct reflect.Value, _ string, names ...string) (StructRules, error) {
	var rules Rules
	var err error
	// TODO: проверить увеличивает ли это производительность
	// if key != "" {
	//	if r, ok := cache[key]; ok {
	//		return r, nil
	//	}
	// }
	structRules := make(StructRules)
	for i := 0; i < rStruct.NumField(); i++ {
		sf := rStruct.Type().Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("validate")
		if tag == "" {
			continue
		}
		nesChecked := false
		fVal := rStruct.Field(i)
		// TODO: подумать об указателях
		fType := fVal.Type()
		switch fType.Kind() { //nolint:exhaustive // есть default
		case reflect.Array, reflect.Slice:
			rules, err = parseTag(fVal.Type().Elem().Kind(), tag)
			structRules[i] = rules
			nesChecked = true
		case reflect.Struct:
			var nested StructRules
			key := ""
			if fType.Name() != "" {
				key = fType.Name() + ":" + fType.PkgPath()
			}
			nested, err = retrieveRules(fVal, key, names...)
			structRules[i] = nested
			nesChecked = true
		// case reflect.Map:
		// case reflect.Bool:
		default:
			rules, err = parseTag(fVal.Kind(), tag)
			structRules[i] = rules
		}
		if tag == "nested" && !nesChecked {
			err = ErrInputNotStruct
		}
		if err != nil {
			return StructRules{}, errors.Wrapf(
				err, "error retrieve rule on `%s`", strings.Join(append(names, sf.Name), "."))
		}
	}
	return structRules, nil
}

func checkStruct(rStruct reflect.Value, rules StructRules, names ...string) ValidationErrors {
	errorSet := make(ValidationErrors, 0)
	for i, ruleSet := range rules {
		sVal := rStruct.Field(i)
		fVal := reflect.Indirect(sVal)
		names := append(names, rStruct.Type().Field(i).Name)
		switch fVal.Kind() { //nolint:exhaustive // есть default
		case reflect.Slice, reflect.Array:
			rules, _ := ruleSet.(Rules)
			for i := 0; i < fVal.Len(); i++ {
				errorSet = append(errorSet, checkValue(fVal.Index(i), rules, append(names, strconv.Itoa(i))...)...)
			}
		case reflect.Struct:
			sRules, _ := ruleSet.(StructRules)
			errorSet = append(errorSet, checkStruct(fVal, sRules, names...)...)
		default:
			rules, _ := ruleSet.(Rules)
			errorSet = append(errorSet, checkValue(fVal, rules, names...)...)
		}
	}

	return errorSet
}

func checkValue(value reflect.Value, rules Rules, names ...string) ValidationErrors {
	errorSet := make(ValidationErrors, 0)
	for _, rule := range rules {
		if err := rule.Check(value); err != nil {
			errorSet = append(errorSet, ValidationError{
				Field: strings.Join(names, "."),
				Err:   err,
			})
		}
	}
	return errorSet
}

func parseTag(kind reflect.Kind, tag string) (Rules, error) {
	tagRules := strings.Split(tag, "|")
	rules := make(Rules, len(tagRules))
	for i, r := range tagRules {
		pos := strings.Index(r, ":")
		if pos < 0 {
			return Rules{}, fmt.Errorf("rule `%s` not found", "")
		}
		ruleID := r[0:pos]
		rule, ok := getRule(ruleID)
		if !ok {
			return Rules{}, fmt.Errorf("rule `%s` not found", ruleID)
		}
		if err := rule.Init(kind, strings.Split(r[pos+1:], ",")); err != nil {
			return Rules{}, err
		}
		rules[i] = rule
	}
	return rules, nil
}

// var cache = make(map[string]StructRules, 0)

/*
var cacheOn = true

func CacheOn() {
	cacheOn = true
}

func CacheOff() {
	cacheOn = false
}
*/
