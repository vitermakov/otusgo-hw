package dto

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

var (
	ErrDateWrongFormat       = errors.New("неверный формат даты начала события, ожидается RFC3339")
	ErrDurationWrongFormat   = errors.New("неверный формат продолжительности события, ожидается 45m, 1h30m")
	ErrNotifyTermWrongFormat = errors.New("неверный формат срока оповещения, ожидается 15m, 3d, 1w")
)

type EventCreate struct {
	Title       string  `json:"title"`
	Date        string  `json:"date"`
	Duration    string  `json:"duration"`    // с единицей измерения.
	Description *string `json:"description"` // опционально.
	NotifyTerm  *string `json:"notifyTerm"`  // опционально, с единицей измерения.
}

// Model возвращает связанную модель model.EventCreate.
func (ec EventCreate) Model() (model.EventCreate, errx.ValidationErrors) {
	var errs errx.ValidationErrors
	input := model.EventCreate{
		Title: ec.Title,
	}
	if date, err := time.Parse(time.RFC3339, ec.Date); err != nil {
		errs.Add(errx.ValidationError{Field: "date", Err: errors.Wrap(ErrDateWrongFormat, err.Error())})
	} else {
		input.Date = date
	}
	if duration, err := time.ParseDuration(ec.Duration); err != nil {
		errs.Add(errx.ValidationError{Field: "duration", Err: errors.Wrap(ErrDurationWrongFormat, err.Error())})
	} else {
		input.Duration = duration
	}
	if ec.Description != nil {
		val := *ec.Description
		input.Description = &val
	}
	if ec.NotifyTerm != nil {
		if notifyTerm, err := time.ParseDuration(*ec.NotifyTerm); err != nil {
			errs.Add(errx.ValidationError{Field: "notifyTerm", Err: errors.Wrap(ErrNotifyTermWrongFormat, err.Error())})
		} else {
			input.NotifyTerm = &notifyTerm
		}
	}
	if errs.Empty() {
		return input, nil
	}
	return model.EventCreate{}, errs
}

type EventUpdate struct {
	Title       *string `json:"title"`
	Date        *string `json:"date"`
	Duration    *string `json:"duration"` // с единицей измерения.
	Description *string `json:"description"`
	NotifyTerm  *string `json:"notifyTerm"` // с единицей измерения.
}

func (eu EventUpdate) Model() (model.EventUpdate, errx.ValidationErrors) {
	var errs errx.ValidationErrors
	input := model.EventUpdate{}
	if eu.Title != nil {
		input.Title = eu.Title
	}
	if eu.Date != nil {
		date, err := time.Parse(time.RFC3339, *eu.Date)
		if err != nil {
			errs.Add(errx.ValidationError{Field: "Date", Err: errors.Wrap(ErrDateWrongFormat, err.Error())})
		} else {
			input.Date = &date
		}
	}
	if eu.Duration != nil {
		if duration, err := time.ParseDuration(*eu.Duration); err != nil {
			errs.Add(errx.ValidationError{Field: "duration", Err: errors.Wrap(ErrDurationWrongFormat, err.Error())})
		} else {
			input.Duration = &duration
		}
	}
	if eu.Description != nil {
		input.Description = eu.Description
	}
	if eu.NotifyTerm != nil {
		if notifyTerm, err := time.ParseDuration(*eu.NotifyTerm); err != nil {
			errs.Add(errx.ValidationError{Field: "notifyTerm", Err: errors.Wrap(ErrNotifyTermWrongFormat, err.Error())})
		} else {
			input.NotifyTerm = &notifyTerm
		}
	}
	if errs.Empty() {
		return input, nil
	}
	return model.EventUpdate{}, errs
}

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Duration    string    `json:"duration"`
	Owner       *User     `json:"owner,omitempty"`
	Description string    `json:"description"`
	NotifyTerm  string    `json:"notifyTerm"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func FromEventModel(item model.Event) Event {
	event := Event{
		ID:          item.ID.String(),
		Title:       item.Title,
		Date:        item.Date,
		Duration:    item.Duration.String(),
		Description: item.Description,
		NotifyTerm:  item.NotifyTerm.String(),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
	if item.Owner != nil {
		user := FromUserModel(*item.Owner)
		event.Owner = &user
	}
	return event
}

func FromEventSlice(items []model.Event) []Event {
	if items == nil {
		return nil
	}
	result := make([]Event, len(items))
	for i, item := range items {
		result[i] = FromEventModel(item)
	}
	return result
}
