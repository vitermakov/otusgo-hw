package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
)

type Event interface {
	Add(context.Context, model.EventCreate) (*model.Event, error)
	Update(context.Context, model.Event, model.EventUpdate) error
	Delete(context.Context, model.Event) error
	GetEventsOnDay(context.Context, model.User, time.Time) ([]model.Event, error)
	GetEventsOnWeek(context.Context, model.User, time.Time) ([]model.Event, error)
	GetEventsOnMonth(context.Context, model.User, time.Time) ([]model.Event, error)
	GetByID(context.Context, uuid.UUID) (*model.Event, error)
}

// User работы с пользователями.
type User interface {
	Add(context.Context, model.UserCreate) (*model.User, error)
	Update(context.Context, model.User, model.UserUpdate) error
	Delete(context.Context, model.User) error
	GetAll(context.Context) ([]model.User, error)
	GetByID(context.Context, uuid.UUID) (*model.User, error)
	GetByEmail(context.Context, string) (*model.User, error)
	// GetCurrent user_id передается в контекстe
	GetCurrent(context.Context) (*model.User, error)
}
