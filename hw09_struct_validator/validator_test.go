package hw09_struct_validator_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09_struct_validator"
	"github.com/vitermakov/otusgo-hw/hw09_struct_validator/validator"
)

func TestValidate(t *testing.T) {
	err := hw09_struct_validator.Validate(nil)
	require.True(t, errors.Is(err, validator.ErrInputStructIsNull))

	err = hw09_struct_validator.Validate(32)
	require.True(t, errors.Is(err, validator.ErrInputNotStruct))
}
