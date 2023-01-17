package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

// Event модель события.
type Event struct {
	ID           uuid.UUID
	Title        string
	Date         time.Time
	Duration     time.Duration
	Owner        *User
	Description  string
	NotifyTerm   time.Duration
	NotifyStatus NotifyStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// EventCreate модель создания события.
type EventCreate struct {
	Title    string
	Date     time.Time
	Duration time.Duration
	OwnerID  uuid.UUID
	// Description описание опционально.
	Description *string
	// NotifyTerm время до начала события для оповещения, опционально.
	NotifyTerm *time.Duration
}

// Validate базовая валидация структуры.
func (ec EventCreate) Validate() error {
	var errs errx.NamedErrors
	if ec.Title == "" {
		errs.Add(errx.NamedError{
			Field: "Title",
			Err:   ErrEventEmptyTitle,
		})
	}
	if ec.Date.IsZero() {
		errs.Add(errx.NamedError{
			Field: "Date",
			Err:   ErrEventZeroDate,
		})
	}
	if ec.Duration <= 0 {
		errs.Add(errx.NamedError{
			Field: "Duration",
			Err:   ErrEventWrongDuration,
		})
	}
	if ec.OwnerID.ID() == 0 {
		errs.Add(errx.NamedError{
			Field: "OwnerID",
			Err:   ErrEventOwnerID,
		})
	}
	if ec.NotifyTerm != nil && *ec.NotifyTerm <= 0 {
		errs.Add(errx.NamedError{
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
	Title        *string
	Date         *time.Time
	Duration     *time.Duration
	Description  *string
	NotifyTerm   *time.Duration
	NotifyStatus *NotifyStatus
}

// Validate базовая валидация структуры.
func (ec EventUpdate) Validate() error {
	var errs errx.NamedErrors
	if ec.Title != nil && *ec.Title == "" {
		errs.Add(errx.NamedError{
			Field: "Title",
			Err:   ErrEventEmptyTitle,
		})
	}
	if ec.Date != nil && ec.Date.IsZero() {
		errs.Add(errx.NamedError{
			Field: "Date",
			Err:   ErrEventZeroDate,
		})
	}
	if ec.Duration != nil && *ec.Duration <= 0 {
		errs.Add(errx.NamedError{
			Field: "Duration",
			Err:   ErrEventWrongDuration,
		})
	}
	if ec.NotifyTerm != nil && *ec.NotifyTerm <= 0 {
		errs.Add(errx.NamedError{
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
	ID        *uuid.UUID
	NotID     *uuid.UUID
	OwnerID   *uuid.UUID
	DateRange *DateRange
	// TacDuration учитывать продолжительность мероприятий.
	TacDuration bool
	// DateLess выбрать события с запланированной датой, меньшей указанной.
	DateLess *time.Time
	// NeedNotifyTerm необходимые к оповещению на указанную дату.
	NeedNotifyTerm *time.Time
}

func EventSearchID(guid string) (EventSearch, error) {
	id, err := uuid.Parse(guid)
	if err != nil {
		return EventSearch{}, err
	}
	return EventSearch{ID: &id}, nil
}
