package dto

import (
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/jsontype"
)

type EventCreate struct {
	Title       jsontype.String  `json:"title"`
	Date        jsontype.String  `json:"date"`
	Duration    jsontype.Int     `json:"duration"`    // в минутах.
	Description *jsontype.String `json:"description"` // опционально.
	NotifyTerm  *jsontype.Int    `json:"notifyTerm"`  // в днях, опционально.
}

// Model возвращает связанную модель если какие-то поля заполены
// в неверном формате, то не возвращаем ошибку, ошибка проверяется в сервисах.
func (ec EventCreate) Model() model.EventCreate {
	input := model.EventCreate{
		Title:    string(ec.Title),
		Duration: time.Duration(ec.Duration) * time.Minute,
	}
	input.Date, _ = time.Parse(time.RFC3339, string(ec.Date))
	if ec.Description != nil {
		val := string(*ec.Description)
		input.Description = &val
	}
	if ec.NotifyTerm != nil {
		val := time.Duration(*ec.NotifyTerm) * time.Hour * 24
		input.NotifyTerm = &val
	}
	return input
}

type EventUpdate struct {
	Title       *jsontype.String `json:"title"`
	Date        *jsontype.String `json:"date"`
	Duration    *jsontype.Int    `json:"duration"` // в минутах.
	Description *jsontype.String `json:"description"`
	NotifyTerm  *jsontype.Int    `json:"notifyTerm"`
}

func (eu EventUpdate) Model() model.EventUpdate {
	input := model.EventUpdate{}
	if eu.Title != nil {
		val := string(*eu.Title)
		input.Title = &val
	}
	if eu.Date != nil {
		// ошибки здесь не проверяем.
		date, _ := time.Parse(time.RFC3339, string(*eu.Date))
		input.Date = &date
	}
	if eu.Duration != nil {
		val := time.Duration(*eu.Duration) * time.Minute
		input.Duration = &val
	}
	if eu.Description != nil {
		val := string(*eu.Description)
		input.Description = &val
	}
	if eu.NotifyTerm != nil {
		val := time.Duration(*eu.NotifyTerm) * time.Hour * 24
		input.NotifyTerm = &val
	}
	return input
}

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Duration    int       `json:"duration"`
	Owner       *User     `json:"owner,omitempty"`
	Description string    `json:"description"`
	NotifyTerm  int       `json:"notifyTerm"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func FromEventModel(item model.Event) Event {
	event := Event{
		ID:          item.ID.String(),
		Title:       item.Title,
		Date:        item.Date,
		Duration:    int(item.Duration.Minutes()),
		Description: item.Description,
		NotifyTerm:  int(item.NotifyTerm.Hours() / 24),
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
