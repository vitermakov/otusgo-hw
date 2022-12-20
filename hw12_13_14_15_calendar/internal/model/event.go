package model

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
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
	NotifyTerm  time.Duration
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// EventCreate модель создания события.
type EventCreate struct {
	Title       string
	Date        time.Time
	Duration    int // в минутах.
	OwnerId     uuid.UUID
	Description *string // опционально.
	NotifyTerm  *int    // в днях, опционально.
}

// Validate базовая валидация структуры
func (ec EventCreate) Validate() error {
	var errs errx.ValidationErrors
	if ec.Title == "" {
		errs.Add(errx.ValidationError{
			Field: "Title",
			Err:   ErrEventEmptyTitle,
		})
	}
	if ec.Date.IsZero() {
		errs.Add(errx.ValidationError{
			Field: "Date",
			Err:   ErrEventZeroDate,
		})
	}
	if ec.Duration <= 0 {
		errs.Add(errx.ValidationError{
			Field: "Duration",
			Err:   ErrEventWrongDuration,
		})
	}
	if ec.OwnerId.ID() == 0 {
		errs.Add(errx.ValidationError{
			Field: "OwnerId",
			Err:   ErrEventOwnerId,
		})
	}
	if ec.NotifyTerm != nil && *ec.NotifyTerm <= 0 {
		errs.Add(errx.ValidationError{
			Field: "NotifyTerm",
			Err:   ErrEventWrongNotifyTerm,
		})
	}
	if errs.Empty() {
		return nil
	}
	return errs
}

// EventUpdate модель обновления события - обновлять можно не все поля.
type EventUpdate struct {
	Title       *string
	Date        *time.Time
	Duration    *int
	Description *string
	NotifyTerm  *int
}

// Validate базовая валидация структуры
func (ec EventUpdate) Validate() error {
	var errs errx.ValidationErrors
	if ec.Title != nil && *ec.Title == "" {
		errs.Add(errx.ValidationError{
			Field: "Title",
			Err:   ErrEventEmptyTitle,
		})
	}
	if ec.Date != nil && ec.Date.IsZero() {
		errs.Add(errx.ValidationError{
			Field: "Date",
			Err:   ErrEventZeroDate,
		})
	}
	if ec.Duration != nil && *ec.Duration <= 0 {
		errs.Add(errx.ValidationError{
			Field: "Duration",
			Err:   ErrEventWrongDuration,
		})
	}
	if ec.NotifyTerm != nil && *ec.NotifyTerm <= 0 {
		errs.Add(errx.ValidationError{
			Field: "NotifyTerm",
			Err:   ErrEventWrongNotifyTerm,
		})
	}
	if errs.Empty() {
		return nil
	}
	return errs
}

// EventSearch модель поиска. Исходя из условия задачи и всех ее аспектов искать события
// необходимо по идентификатору и промежутку дат (с учетом и без учета продолжительности).
type EventSearch struct {
	ID          *uuid.UUID
	NotID       *uuid.UUID
	OwnerID     *uuid.UUID
	DateRange   *DateRange
	TacDuration bool // учитывать продолжительность мероприятий.
}

func EventSearchID(guid string) (EventSearch, error) {
	id, err := uuid.Parse(guid)
	if err != nil {
		return EventSearch{}, err
	}
	return EventSearch{ID: &id}, nil
}
