package dto

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
)

type UserCreate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (uc UserCreate) Model() model.UserCreate {
	return model.UserCreate{
		Name:  uc.Name,
		Email: uc.Email,
	}
}

type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (uu UserUpdate) Model() model.UserUpdate {
	input := model.UserUpdate{}
	if uu.Name != nil {
		input.Name = uu.Name
	}
	if uu.Email != nil {
		input.Email = uu.Email
	}
	return input
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func FromUserModel(item model.User) User {
	return User{
		ID:    item.ID.String(),
		Name:  item.Name,
		Email: item.Email,
	}
}

func FromUserSlice(items []model.User) []User {
	if items == nil {
		return nil
	}
	result := make([]User, len(items))
	for i, item := range items {
		result[i] = FromUserModel(item)
	}
	return result
}
