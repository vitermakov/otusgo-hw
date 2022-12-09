package hw09structvalidator

import (
	"github.com/vitermakov/otusgo-hw/hw09structvalidator/validator"
)

func Validate(v interface{}) error {
	return validator.ValidateStruct(v)
}
