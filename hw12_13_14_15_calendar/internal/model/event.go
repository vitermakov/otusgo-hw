package model

import (
	"time"

	"github.com/google/uuid"
)

// Event модель события
type Event struct {
	ID          uuid.UUID
	Title       string
	Date        time.Time
	Duration    time.Duration
	Owner       *User
	Description string
	NotifyTerm  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// EventCreate модель создания события.
type EventCreate struct {
	Title       string
	Date        time.Time
	Duration    time.Duration
	OwnerId     *uuid.UUID
	Description *string // опционально
	NotifyTerm  *int    // опционально
}

// EventUpdate модель обновления события - обновлять можно не все поля.
type EventUpdate struct {
	Title       *string
	Date        *time.Time
	Duration    *time.Duration
	Description *string
	NotifyTerm  *int
}

// EventSearch модель поиска. Исходя из условия задачи и всех ее аспектов искать события
// необходимо по идентификатору и промежутку дат.
type EventSearch struct {
	ID        *uuid.UUID
	NotID     *uuid.UUID
	DateRange *DateRange
}

func EventSearchID(guid string) (EventSearch, error) {
	id, err := uuid.Parse(guid)
	if err != nil {
		return EventSearch{}, err
	}
	return EventSearch{ID: &id}, nil
}
