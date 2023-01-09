package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
)

// EventCRUD сервис управления пользователями. Во всех случаях считаем,
// что добавлять, удалять... можно только "свои" события.
type EventCRUD interface {
	Add(context.Context, model.EventCreate) (*model.Event, error)
	Update(context.Context, model.Event, model.EventUpdate) error
	Delete(context.Context, model.Event) error
	GetUserEventsOn(context.Context, time.Time, model.RangeKind) ([]model.Event, error)
	GetEvents(context.Context, model.EventSearch) ([]model.Event, error)
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
	// GetCurrent user_id передается в контексте.
	GetCurrent(context.Context) (*model.User, error)
}

// EventNotify сервис управления оповещениями.
type EventNotify interface {
	GetNotifications(context.Context) ([]model.Notification, error)
	MarkEventNotified(context.Context, uuid.UUID) error
}

// EventClean удаление устаревших объектов календаря.
type EventClean interface {
	CleanupOldEvents(context.Context, time.Duration) (int64, error)
}
