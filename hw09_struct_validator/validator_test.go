package hw09structvalidator_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09structvalidator"
	"github.com/vitermakov/otusgo-hw/hw09structvalidator/validator"
)

func TestValidate(t *testing.T) {
	err := hw09structvalidator.Validate(nil)
	require.True(t, errors.Is(err, validator.ErrInputStructIsNull))

	err = hw09structvalidator.Validate(32)
	require.True(t, errors.Is(err, validator.ErrInputNotStruct))
}
