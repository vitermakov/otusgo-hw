package repository

import (
	"context"
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
)

// Event репозиторий для управления событиями.
// Этот интерфейс не реализует никакой бизнес логики, его задача - взаимодействовать с хранилищем.
// Делаем методы максимально общими.
type Event interface {
	Add(context.Context, model.EventCreate) (*model.Event, error)
	Update(context.Context, model.EventUpdate, model.EventSearch) error
	Delete(context.Context, model.EventSearch) error
	// GetList не учитываем пагинацию и сортировку.
	GetList(context.Context, model.EventSearch) ([]model.Event, error)
	BlockEvents4Notify(context.Context, time.Time) ([]model.Event, error)
}

// User репозиторий для управления пользователями.
type User interface {
	Add(context.Context, model.UserCreate) (*model.User, error)
	Update(context.Context, model.UserUpdate, model.UserSearch) error
	Delete(context.Context, model.UserSearch) error
	// GetList не учитываем пагинацию и сортировку.
	GetList(context.Context, model.UserSearch) ([]model.User, error)
}
