package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type UserService struct {
	repo repository.User
	log  logger.Logger
}

func (us UserService) validateAdd(ctx context.Context, input model.UserCreate) error {
	users, err := us.repo.GetList(ctx, model.UserSearch{
		Email: &input.Email,
	})
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return model.ErrUserDuplicateEmail
	}
	return nil
}
func (us UserService) Add(ctx context.Context, input model.UserCreate) (*model.User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := us.validateAdd(ctx, input); err != nil {
		return nil, err
	}
	return us.repo.Add(ctx, input)
}

func (us UserService) validateUpdate(ctx context.Context, user model.User, input model.UserUpdate) error {
	if input.Email != nil {
		users, err := us.repo.GetList(ctx, model.UserSearch{
			Email: input.Email,
		})
		if err != nil {
			return err
		}
		if len(users) > 0 && users[0].ID.String() != user.ID.String() {
			return model.ErrUserDuplicateEmail
		}
	}

	return nil
}

func (us UserService) Update(ctx context.Context, user model.User, input model.UserUpdate) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if err := us.validateUpdate(ctx, user, input); err != nil {
		return err
	}
	return us.repo.Update(ctx, input, model.UserSearch{ID: &user.ID})
}

func (us UserService) Delete(ctx context.Context, user model.User) error {
	if err := us.repo.Delete(ctx, model.UserSearch{ID: &user.ID}); err != nil {
		return err
	}
	return nil
}

func (us UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return us.repo.GetList(ctx, model.UserSearch{})
}

func (us UserService) GetById(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	return us.getOne(ctx, model.UserSearch{ID: &userId})
}

func (us UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return us.getOne(ctx, model.UserSearch{Email: &email})
}

func (us UserService) getOne(ctx context.Context, search model.UserSearch) (*model.User, error) {
	users, err := us.repo.GetList(ctx, search)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, model.ErrUserNotFound
	}
	return &users[0], nil
}

func (us UserService) GetCurrent(ctx context.Context) (*model.User, error) {
	rawId, ok := ctx.Value("user_id").(string)
	if !ok || rawId == "" {
		return nil, model.ErrUserEmptyID
	}
	userId, err := uuid.Parse(rawId)
	if err != nil {
		return nil, err
	}
	return us.GetById(ctx, userId)
}

func NewUserService(repo repository.User, log logger.Logger) User {
	return &UserService{
		repo: repo,
		log:  log,
	}
}
