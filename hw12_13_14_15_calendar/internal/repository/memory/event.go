package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
)

type EventRepo struct {
	mu     sync.RWMutex
	events []model.Event
}

func NewEventRepo() repository.Event {
	return &EventRepo{}
}

func (er *EventRepo) Add(ctx context.Context, input model.EventCreate) (*model.Event, error) {
	event := model.Event{
		ID:        uuid.New(),
		Title:     input.Title,
		Date:      input.Date,
		Duration:  input.Duration,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if input.OwnerID.ID() > 0 {
		event.Owner = &model.User{ID: input.OwnerID}
	}
	if input.Description != nil {
		event.Description = *input.Description
	}
	if input.NotifyTerm != nil {
		event.NotifyTerm = *input.NotifyTerm
	}
	er.mu.Lock()
	er.events = append(er.events, event)
	er.mu.Unlock()

	return &event, nil
}

func (er *EventRepo) Update(ctx context.Context, input model.EventUpdate, search model.EventSearch) error {
	er.mu.Lock()
	for i, event := range er.events {
		if !er.matchSearch(event, search) {
			continue
		}
		if input.Title != nil {
			event.Title = *input.Title
		}
		if input.Date != nil {
			event.Date = *input.Date
		}
		if input.Duration != nil {
			event.Duration = *input.Duration
		}
		if input.Description != nil {
			event.Description = *input.Description
		}
		if input.NotifyTerm != nil {
			event.NotifyTerm = *input.NotifyTerm
		}
		event.UpdatedAt = time.Now()
		er.events[i] = event
	}
	er.mu.Unlock()
	return nil
}

func (er *EventRepo) Delete(ctx context.Context, search model.EventSearch) error {
	er.mu.Lock()
	result := make([]model.Event, 0)
	for _, event := range er.events {
		if !er.matchSearch(event, search) {
			result = append(result, event)
		}
	}
	er.events = result
	er.mu.Unlock()
	return nil
}

// GetList не учитываем пагинацию, сортировку.
func (er *EventRepo) GetList(ctx context.Context, search model.EventSearch) ([]model.Event, error) {
	var filtered []model.Event
	er.mu.RLock()
	for _, event := range er.events {
		if er.matchSearch(event, search) {
			filtered = append(filtered, event)
		}
	}
	er.mu.RUnlock()

	return filtered, nil
}

func (er *EventRepo) matchSearch(event model.Event, search model.EventSearch) bool {
	if search.ID != nil {
		if strings.Compare(event.ID.String(), search.ID.String()) != 0 {
			return false
		}
	}
	if search.NotID != nil {
		if strings.Compare(event.ID.String(), search.NotID.String()) == 0 {
			return false
		}
	}
	if search.OwnerID != nil {
		if strings.Compare(event.Owner.ID.String(), search.OwnerID.String()) != 0 {
			return false
		}
	}
	if search.DateRange != nil {
		evStart := event.Date
		evEnd := event.Date
		if search.TacDuration {
			evEnd = evEnd.Add(event.Duration)
		}
		if !(evEnd.After(search.DateRange.GetFrom()) &&
			evStart.Before(search.DateRange.GetTo())) {
			return false
		}
	}
	if search.DateLess != nil && event.Date.After(*search.DateLess) {
		return false
	}
	if search.NeedNotifyTerm != nil {
		if event.NotifyTerm.Seconds() == 0 {

		}
		if event.Date.After(*search.DateLess) {
			return false
		}
	}
	return true
}

func (er *EventRepo) BlockEvents4Notify(ctx context.Context, now time.Time) ([]model.Event, error) {
	return nil, nil
}
