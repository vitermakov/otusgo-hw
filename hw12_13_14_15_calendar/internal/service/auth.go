package service

import (
	"context"
	"errors"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

type AuthService struct {
	user User
}

func (as AuthService) Authorize(ctx context.Context, email string) (*servers.AuthUser, error) {
	user, err := as.user.GetByEmail(ctx, email)
	if err != nil {
		// пользователь не найден?
		nfErr := errx.NotFound{}
		if errors.As(err, &nfErr) {
			return nil, nil
		}
		return nil, err
	}
	return &servers.AuthUser{
		ID:    user.ID.String(),
		Login: user.Email,
		Name:  user.Name,
	}, nil
}

func NewAuthService(user User) servers.AuthService {
	return &AuthService{user}
}
