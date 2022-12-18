package model

import (
	"github.com/google/uuid"
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

// UserUpdate модель изменения пользователя.
type UserUpdate struct {
	Name  *string
	Email *string
}

// UserSearch модель поиска пользователя.
type UserSearch struct {
	ID *uuid.UUID
}
