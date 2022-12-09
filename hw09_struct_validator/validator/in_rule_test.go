package validator_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09structvalidator/validator"
)

func TestInRule(t *testing.T) {
	allowedValues := []string{"new", "running", "done", "error"}
	wrongValue := reflect.ValueOf("complete")
	rightValue := reflect.ValueOf("done")
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
			name: "wrong type",
			kind: reflect.Map,
			args: allowedValues,
			assertInitErr: func(t *testing.T, err error) {
				t.Helper()
				require.True(t, errors.Is(err, validator.ErrSupportArgType))
			},
		}, {
			name: "wrong value",
			kind: reflect.String,
			args: allowedValues,
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
			args: allowedValues,
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
			rule := validator.NewInRule()
			err := rule.Init(tc.kind, tc.args)
			tc.assertInitErr(t, err)

			if tc.checkValue != nil {
				err := rule.Check(*tc.checkValue)
				tc.assertCheckErr(t, err)
			}
		})
	}
}
