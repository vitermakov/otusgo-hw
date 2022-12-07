package validator

import (
	"fmt"
	"reflect"
)

// Rule интерфейс правила валидации
type Rule interface {
	// Инициализация значения идет из tag-параметров, поэтому как исходное значение передается
	// массив строк, дальнейшая инициализация которыми конкретного правила идет индивидуально.
	Init(reflect.Kind, []string) error

	// Проверка значения
	Check(reflect.Value) error
}

// RuleCreate конструктор правила
type RuleCreate func() Rule

// Registry регистр правил.
// - в качестве ключа используется строковый идентификатор правила валидации
// - в качестве значения - конструктор правила
type Registry map[string]RuleCreate

// имеющиеся на данный момент правила
var registry = Registry{
	"min":    NewMinRule,
	"max":    NewMaxRule,
	"in":     NewInRule,
	"len":    NewLenRule,
	"regexp": NewReRule,
}

// поиск конструктора правил.
// возвращает созданный объект правила или nil и флаг существования указанного правила
func getRule(key string) (Rule, bool) {
	create, ok := registry[key]
	if !ok {
		return nil, false
	}
	return create(), true
}

type Invalid struct {
	Code string
	Err  error
}

func (i Invalid) Error() string {
	return fmt.Sprintf("(%s) %s", i.Code, i.Err.Error())
}
