package validator

import (
	"fmt"
	"reflect"
)

// Rule интерфейс правила валидации.
type rule interface {
	// Проверка значения
	Check(reflect.Value) error
}

// RuleCreate конструктор правила.
// Инициализация значения идет из tag-параметров, поэтому как исходное значение передается
// массив строк, дальнейшая инициализация которыми конкретного правила идет индивидуально.
type ruleCreate func(kind reflect.Kind, args []string) (rule, error)

// ruleRegistry регистр правил
// - в качестве ключа используется строковый идентификатор правила валидации
// - в качестве значения - конструктор правила.
type ruleRegistry map[string]ruleCreate

// имеющиеся на данный момент правила.
var registry = ruleRegistry{
	"min":    createMinRule,
	"max":    createMaxRule,
	"in":     createInRule,
	"len":    createLenRule,
	"regexp": createReRule,
}

// абстрактная фабрика правил.
func GetRuleFactory(key string, kind reflect.Kind, args []string) (rule, error) {
	create, ok := registry[key]
	if !ok {
		return nil, fmt.Errorf("rule '%s' not exists", key)
	}
	return create(kind, args)
}

// ошибка, возвращаемая при ошибке валидации конкретного обработчика.
type Invalid struct {
	Code string
	Err  error
}

func (i Invalid) Error() string {
	return fmt.Sprintf("(%s) %s", i.Code, i.Err.Error())
}
