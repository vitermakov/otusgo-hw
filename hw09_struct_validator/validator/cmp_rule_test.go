package validator_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09structvalidator/validator"
)

func TestMinRule(t *testing.T) {
	v4 := reflect.ValueOf(4)
	v10 := reflect.ValueOf(10)
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
			kind: reflect.Int,
			args: []string{},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name: "wrong argrs count (>1)",
			kind: reflect.Int,
			args: []string{"4", "3"},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name: "arg not number",
			kind: reflect.Int,
			args: []string{"xxx"},
			assertInitErr: func(t *testing.T, err error) {
				_, ok := err.(*strconv.NumError)
				require.True(t, ok)
			},
		}, {
			name: "wrong type",
			kind: reflect.Map,
			args: []string{"4"},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrSupportArgType))
			},
		}, {
			name: "wrong value",
			kind: reflect.Float32,
			args: []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkValue: &v4,
			assertCheckErr: func(t *testing.T, err error) {
				_, ok := err.(validator.Invalid)
				require.True(t, ok)
			},
		}, {
			name: "ok",
			kind: reflect.Float32,
			args: []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkValue: &v10,
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
			rule := validator.NewMinRule()
			err := rule.Init(tc.kind, tc.args)
			tc.assertInitErr(t, err)

			if tc.checkValue != nil {
				err := rule.Check(*tc.checkValue)
				tc.assertCheckErr(t, err)
			}
		})
	}
}

func TestMaxRule(t *testing.T) {
	v4 := reflect.ValueOf(4)
	v10 := reflect.ValueOf(10)
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
			kind: reflect.Int,
			args: []string{},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name: "wrong argrs count (>1)",
			kind: reflect.Int,
			args: []string{"4", "3"},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrWrongArgsList))
			},
		}, {
			name: "arg not number",
			kind: reflect.Int,
			args: []string{"xxx"},
			assertInitErr: func(t *testing.T, err error) {
				_, ok := err.(*strconv.NumError)
				require.True(t, ok)
			},
		}, {
			name: "wrong type",
			kind: reflect.Map,
			args: []string{"4"},
			assertInitErr: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, validator.ErrSupportArgType))
			},
		}, {
			name: "wrong value",
			kind: reflect.Float32,
			args: []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkValue: &v10,
			assertCheckErr: func(t *testing.T, err error) {
				_, ok := err.(validator.Invalid)
				require.True(t, ok)
			},
		}, {
			name: "ok",
			kind: reflect.Float32,
			args: []string{"7.0"},
			assertInitErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkValue: &v4,
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
			rule := validator.NewMaxRule()
			err := rule.Init(tc.kind, tc.args)
			tc.assertInitErr(t, err)

			if tc.checkValue != nil {
				err := rule.Check(*tc.checkValue)
				tc.assertCheckErr(t, err)
			}
		})
	}
}
