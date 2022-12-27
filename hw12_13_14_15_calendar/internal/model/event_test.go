package model

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

func TestEventCreateValidate(t *testing.T) {
	wrongNotifyTerm := -1 * time.Hour
	testCases := []struct {
		name     string
		input    EventCreate
		expected errx.ValidationErrors
	}{
		{
			name: "bad event create",
			input: EventCreate{
				NotifyTerm: &wrongNotifyTerm,
			},
			expected: []errx.ValidationError{
				{
					Field: "Title",
					Err:   ErrEventEmptyTitle,
				}, {
					Field: "Date",
					Err:   ErrEventZeroDate,
				}, {
					Field: "Duration",
					Err:   ErrEventWrongDuration,
				}, {
					Field: "OwnerID",
					Err:   ErrEventOwnerID,
				}, {
					Field: "NotifyTerm",
					Err:   ErrEventWrongNotifyTerm,
				},
			},
		}, {
			name: "bad event create",
			input: EventCreate{
				Duration: -3,
			},
			expected: []errx.ValidationError{
				{
					Field: "Title",
					Err:   ErrEventEmptyTitle,
				}, {
					Field: "Date",
					Err:   ErrEventZeroDate,
				}, {
					Field: "Duration",
					Err:   ErrEventWrongDuration,
				}, {
					Field: "OwnerID",
					Err:   ErrEventOwnerID,
				},
			},
		}, {
			name: "ok event create",
			input: EventCreate{
				Title:    "Test Event",
				Date:     time.Now(),
				Duration: 10,
				OwnerID:  uuid.New(),
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

func TestEventUpdateValidate(t *testing.T) {
	var (
		emptyTitle      string
		okTitle         = "okEvent"
		zeroDate        time.Time
		okDate          = time.Now()
		wrongDuration   = -4
		okDuration      = 4
		wrongNotifyTerm = -10
		okNotifyTerm    = 20
	)
	testCases := []struct {
		name     string
		input    EventUpdate
		expected errx.ValidationErrors
	}{
		{
			name:     "ok empty event update",
			input:    EventUpdate{},
			expected: nil,
		}, {
			name: "bad event update",
			input: EventUpdate{
				Title:      &emptyTitle,
				Date:       &zeroDate,
				Duration:   &wrongDuration,
				NotifyTerm: &wrongNotifyTerm,
			},
			expected: []errx.ValidationError{
				{
					Field: "Title",
					Err:   ErrEventEmptyTitle,
				}, {
					Field: "Date",
					Err:   ErrEventZeroDate,
				}, {
					Field: "Duration",
					Err:   ErrEventWrongDuration,
				}, {
					Field: "NotifyTerm",
					Err:   ErrEventWrongNotifyTerm,
				},
			},
		}, {
			name: "ok event update",
			input: EventUpdate{
				Title:      &okTitle,
				Date:       &okDate,
				Duration:   &okDuration,
				NotifyTerm: &okNotifyTerm,
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
