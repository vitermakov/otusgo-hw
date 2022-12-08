package validator_test

import (
	"errors"
	"reflect"
	"regexp/syntax"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09_struct_validator/validator"
)

func TestReRule(t *testing.T) {
	mailRe := "\\w+@\\w+\\.\\w+$"
	wrongValue := reflect.ValueOf("mail@yandex.")
	rightValue := reflect.ValueOf("mail@yandex.ru")
	testCases := []struct {
		name           string
		kind           reflect.Kind
		args           []string
		assertInitErr  func(*testing.T, error)
		checkValue     *reflect.Value
		assertCheckErr func(*testing.T, error)
	}{
		{
			name: "wrong argrs count",
			kind: reflect.String,
			args: []string{},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name: "arg re syntax error",
			kind: reflect.String,
			args: []string{"(()"},
			assertInitErr: func(t *testing.T, err error) {
				_, ok := err.(*syntax.Error)
				require.True(t, ok)
			},
		}, {
			name: "wrong type",
			kind: reflect.Int,
			args: []string{mailRe},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrSupportArgType))
			},
		}, {
			name: "wrong value",
			kind: reflect.String,
			args: []string{mailRe},
			assertInitErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkValue: &wrongValue,
			assertCheckErr: func(t *testing.T, err error) {
				_, ok := err.(validator.Invalid)
				require.True(t, ok)
			},
		}, {
			name: "ok",
			kind: reflect.String,
			args: []string{mailRe},
			assertInitErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkValue: &rightValue,
			assertCheckErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			// Place your code here.
			rule := validator.NewReRule()
			err := rule.Init(tc.kind, tc.args)
			tc.assertInitErr(t, err)

			if tc.checkValue != nil {
				err := rule.Check(*tc.checkValue)
				tc.assertCheckErr(t, err)
			}
		})
	}
}
