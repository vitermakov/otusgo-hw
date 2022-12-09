package validator_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09_struct_validator/validator"
)

func TestMinRule(t *testing.T) {
	v4 := reflect.ValueOf(4)
	v10 := reflect.ValueOf(10)
	testCases := []struct {
		name           string
		rules          []string
		kind           reflect.Kind
		args           []string
		assertInitErr  func(*testing.T, error)
		checkValue     *reflect.Value
		assertCheckErr func(*testing.T, error)
	}{
		{
			name:  "wrong argrs count (0)",
			rules: []string{"min", "max"},
			kind:  reflect.Int,
			args:  []string{},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name:  "wrong argrs count (>1)",
			rules: []string{"min", "max"},
			kind:  reflect.Int,
			args:  []string{"4", "3"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name:  "arg not number",
			rules: []string{"min", "max"},
			kind:  reflect.Int,
			args:  []string{"xxx"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				var ne *strconv.NumError
				require.True(t, errors.As(err, &ne))
			},
		}, {
			name:  "wrong type",
			rules: []string{"min", "max"},
			kind:  reflect.Map,
			args:  []string{"4"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, validator.ErrSupportArgType))
			},
		}, {
			name:  "wrong value",
			rules: []string{"min"},
			kind:  reflect.Float32,
			args:  []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
			checkValue: &v4,
			assertCheckErr: func(t *testing.T, err error) {
				t.Helper()
				var ne validator.Invalid
				require.True(t, errors.As(err, &ne))
			},
		}, {
			name:  "ok",
			rules: []string{"min"},
			kind:  reflect.Float32,
			args:  []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
			checkValue: &v10,
			assertCheckErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		}, {
			name:  "wrong value",
			rules: []string{"max"},
			kind:  reflect.Float32,
			args:  []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
			checkValue: &v10,
			assertCheckErr: func(t *testing.T, err error) {
				t.Helper()
				var ne validator.Invalid
				require.True(t, errors.As(err, &ne))
			},
		}, {
			name:  "ok",
			rules: []string{"max"},
			kind:  reflect.Float32,
			args:  []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
			checkValue: &v4,
			assertCheckErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		for _, rc := range tc.rules {
			rc := rc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				rule, err := validator.GetRuleFactory(rc, tc.kind, tc.args)
				tc.assertInitErr(t, err)

				if tc.checkValue != nil {
					err := rule.Check(*tc.checkValue)
					tc.assertCheckErr(t, err)
				}
			})
		}
	}
}
