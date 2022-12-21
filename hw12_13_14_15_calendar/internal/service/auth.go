package service

import (
	"context"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest"
)

type AuthService struct {
	user User
}

func (as AuthService) Authorize(ctx context.Context, email string) (*rest.AuthUser, error) {
	user, err := as.user.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &rest.AuthUser{
		ID:    user.ID.String(),
		Login: user.Email,
		Name:  user.Name,
	}, nil
}

func NewAuthService(user User) rest.AuthService {
	return &AuthService{user}
}
