package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
)

type UserRepo struct {
	mu    sync.RWMutex
	users []model.User
}

func NewUserRepo() repository.User {
	return &UserRepo{}
}

func (ur *UserRepo) Add(ctx context.Context, input model.UserCreate) (*model.User, error) {
	user := model.User{
		ID:    uuid.New(),
		Name:  input.Name,
		Email: input.Email,
	}
	ur.mu.Lock()
	ur.users = append(ur.users, user)
	ur.mu.Unlock()

	return &user, nil
}

func (ur *UserRepo) Update(ctx context.Context, input model.UserUpdate, search model.UserSearch) error {
	ur.mu.Lock()
	for i, user := range ur.users {
		if !ur.matchSearch(user, search) {
			continue
		}
		if input.Name != nil {
			user.Name = *input.Name
		}
		if input.Email != nil {
			user.Email = *input.Email
		}
		ur.users[i] = user
	}
	ur.mu.Unlock()
	return nil
}

func (ur *UserRepo) Delete(ctx context.Context, search model.UserSearch) error {
	ur.mu.Lock()
	result := make([]model.User, 0)
	for _, user := range ur.users {
		if !ur.matchSearch(user, search) {
			result = append(result, user)
		}
	}
	ur.users = result
	ur.mu.Unlock()
	return nil
}

// GetList не учитываем пагинацию, сортировку.
func (ur *UserRepo) GetList(ctx context.Context, search model.UserSearch) ([]model.User, error) {
	var users, filtered []model.User
	ur.mu.RLock()
	users = ur.users
	ur.mu.RUnlock()
	for _, user := range users {
		if ur.matchSearch(user, search) {
			filtered = append(filtered, user)
		}
	}
	return filtered, nil
}

func (ur *UserRepo) matchSearch(user model.User, search model.UserSearch) bool {
	if search.ID != nil {
		if strings.Compare(user.ID.String(), search.ID.String()) != 0 {
			return false
		}
	}
	if search.Email != nil {
		if strings.Compare(user.Email, *search.Email) != 0 {
			return false
		}
	}
	return true
}
