package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

func TestUserCreateValidate(t *testing.T) {
	testCases := []struct {
		name     string
		input    UserCreate
		expected errx.ValidationErrors
	}{
		{
			name:  "bad user create",
			input: UserCreate{},
			expected: []errx.ValidationError{
				{
					Field: "Name",
					Err:   ErrUserEmptyName,
				}, {
					Field: "Email",
					Err:   ErrUserWrongEmail,
				},
			},
		}, {
			name: "bad user add mail",
			input: UserCreate{
				Name:  "Test User",
				Email: "dfvs%dvsdv",
			},
			expected: []errx.ValidationError{
				{
					Field: "Email",
					Err:   ErrUserWrongEmail,
				},
			},
		}, {
			name: "ok user create",
			input: UserCreate{
				Name:  "Test User",
				Email: "test@test.ru",
			},
			expected: nil,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.input.Validate()
			if tc.expected == nil {
				require.NoError(t, err)
			} else {
				var validationError errx.ValidationErrors
				require.True(t, errors.As(err, &validationError))
				require.ElementsMatch(t, validationError, tc.expected)
			}
		})
	}
}

func TestUserUpdateValidate(t *testing.T) {
	var (
		emptyName  string
		okName     = "okName"
		emptyEmail string
		wrongEmail = "erferf#2dfv"
		okEmail    = "test@test.ru"
	)
	testCases := []struct {
		name     string
		input    UserUpdate
		expected errx.ValidationErrors
	}{
		{
			name:     "ok empty user update",
			input:    UserUpdate{},
			expected: nil,
		}, {
			name: "bad user update",
			input: UserUpdate{
				Name:  &emptyName,
				Email: &emptyEmail,
			},
			expected: []errx.ValidationError{
				{
					Field: "Name",
					Err:   ErrUserEmptyName,
				}, {
					Field: "Email",
					Err:   ErrUserWrongEmail,
				},
			},
		}, {
			name: "bad user update mail",
			input: UserUpdate{
				Name:  &okName,
				Email: &wrongEmail,
			},
			expected: []errx.ValidationError{
				{
					Field: "Email",
					Err:   ErrUserWrongEmail,
				},
			},
		}, {
			name: "ok user update",
			input: UserUpdate{
				Name:  &okName,
				Email: &okEmail,
			},
			expected: nil,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.input.Validate()
			if tc.expected == nil {
				require.NoError(t, err)
			} else {
				var validationError errx.ValidationErrors
				require.True(t, errors.As(err, &validationError))
				require.ElementsMatch(t, validationError, tc.expected)
			}
		})
	}
}
