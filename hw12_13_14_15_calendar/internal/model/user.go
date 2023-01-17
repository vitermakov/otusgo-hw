package model

import (
	"net/mail"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

// User модель пользователя.
type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

// UserCreate модель создания пользователя.
type UserCreate struct {
	Name  string
	Email string
}

// Validate базовая валидация структуры.
func (uc UserCreate) Validate() error {
	var errs errx.NamedErrors
	if uc.Name == "" {
		errs.Add(errx.NamedError{
			Field: "Name",
			Err:   ErrUserEmptyName,
		})
	}
	_, err := mail.ParseAddress(uc.Email)
	if err != nil {
		errs.Add(errx.NamedError{
			Field: "Email",
			Err:   ErrUserWrongEmail,
		})
	}
	if errs.Empty() {
		return nil
	}
	return errs
}

// UserUpdate модель изменения пользователя.
type UserUpdate struct {
	Name  *string
	Email *string
}

// Validate базовая валидация структуры.
func (uu UserUpdate) Validate() error {
	var errs errx.NamedErrors
	if uu.Name != nil && *uu.Name == "" {
		errs.Add(errx.NamedError{
			Field: "Name",
			Err:   ErrUserEmptyName,
		})
	}
	if uu.Email != nil {
		_, err := mail.ParseAddress(*uu.Email)
		if err != nil {
			errs.Add(errx.NamedError{
				Field: "Email",
				Err:   ErrUserWrongEmail,
			})
		}
	}
	if errs.Empty() {
		return nil
	}
	return errs
}

// UserSearch модель поиска пользователя.
type UserSearch struct {
	ID    *uuid.UUID
	Email *string
}
