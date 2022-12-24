package dto

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/jsontype"
)

type UserCreate struct {
	Name  jsontype.String `json:"name"`
	Email jsontype.String `json:"email"`
}

func (uc UserCreate) Model() model.UserCreate {
	return model.UserCreate{
		Name:  string(uc.Name),
		Email: string(uc.Email),
	}
}

type UserUpdate struct {
	Name  *jsontype.String `json:"name"`
	Email *jsontype.String `json:"email"`
}

func (uu UserUpdate) Model() model.UserUpdate {
	input := model.UserUpdate{}
	if uu.Name != nil {
		val := string(*uu.Name)
		input.Name = &val
	}
	if uu.Email != nil {
		val := string(*uu.Email)
		input.Email = &val
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
