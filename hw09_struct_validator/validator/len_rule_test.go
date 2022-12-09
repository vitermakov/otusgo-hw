package validator_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09structvalidator/validator"
)

func TestLenRule(t *testing.T) {
	wrongValue := reflect.ValueOf("94506")
	rightValue := reflect.ValueOf("3928194506")
	testCases := []struct {
		name           string
		kind           reflect.Kind
		args           []string
		assertInitErr  func(*testing.T, error)
		checkValue     *reflect.Value
		assertCheckErr func(*testing.T, error)
	}{
		{
			name: "wrong argrs count (0)",
			kind: reflect.String,
			args: []string{},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name: "wrong argrs count (>1)",
			kind: reflect.String,
			args: []string{"4", "3"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name: "arg not int",
			kind: reflect.String,
			args: []string{"xxx"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				var ne *strconv.NumError
				require.True(t, errors.As(err, &ne))
			},
		}, {
			name: "arg not positive",
			kind: reflect.String,
			args: []string{"-10"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorContains(t, err, "len must be positive")
			},
		}, {
			name: "wrong type",
			kind: reflect.Int,
			args: []string{"4"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, validator.ErrSupportArgType))
			},
		}, {
			name: "wrong value",
			kind: reflect.String,
			args: []string{"10"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
			checkValue: &wrongValue,
			assertCheckErr: func(t *testing.T, err error) {
				t.Helper()
				var ne validator.Invalid
				require.True(t, errors.As(err, &ne))
			},
		}, {
			name: "ok",
			kind: reflect.String,
			args: []string{"10"},
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
			checkValue: &rightValue,
			assertCheckErr: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			// Place your code here.
			rule := validator.NewLenRule()
			err := rule.Init(tc.kind, tc.args)
			tc.assertInitErr(t, err)

			if tc.checkValue != nil {
				err := rule.Check(*tc.checkValue)
				tc.assertCheckErr(t, err)
			}
		})
	}
}
