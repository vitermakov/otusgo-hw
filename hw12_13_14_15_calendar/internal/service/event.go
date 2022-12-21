package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type EventService struct {
	repo repository.Event
	log  logger.Logger
	user User
}

func (es EventService) validateAdd(ctx context.Context, input model.EventCreate) error {
	user, err := es.user.GetByID(ctx, input.OwnerID)
	if err != nil {
		return err
	}
	if user == nil {
		return model.ErrEventOwnerExists
	}
	events, err := es.repo.GetList(ctx, model.EventSearch{
		OwnerID: &input.OwnerID,
		DateRange: &model.DateRange{
			DateStart: input.Date,
			Duration:  time.Duration(input.Duration) * time.Minute,
		},
		TacDuration: true,
	})
	if err != nil {
		return err
	}
	if len(events) > 0 {
		return model.ErrEventDateBusy
	}
	return nil
}

func (es EventService) Add(ctx context.Context, input model.EventCreate) (*model.Event, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := es.validateAdd(ctx, input); err != nil {
		return nil, err
	}
	return es.repo.Add(ctx, input)
}

func (es EventService) validateUpdate(ctx context.Context, event model.Event, input model.EventUpdate) error {
	if input.Date != nil || input.Duration != nil {
		dateRgn := model.DateRange{
			DateStart: event.Date,
			Duration:  event.Duration,
		}
		if input.Date != nil {
			dateRgn.DateStart = *input.Date
		}
		if input.Duration != nil {
			dateRgn.Duration = time.Duration(*input.Duration) * time.Minute
		}
		search := model.EventSearch{
			OwnerID:     &event.Owner.ID,
			NotID:       &event.ID,
			DateRange:   &dateRgn,
			TacDuration: true,
		}
		events, err := es.repo.GetList(ctx, search)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			return model.ErrEventDateBusy
		}
	}
	return nil
}

func (es EventService) Update(ctx context.Context, event model.Event, input model.EventUpdate) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if err := es.validateUpdate(ctx, event, input); err != nil {
		return err
	}
	return es.repo.Update(ctx, input, model.EventSearch{ID: &event.ID})
}

func (es EventService) GetEventsOnDay(ctx context.Context, user model.User, date time.Time) ([]model.Event, error) {
	dateRgn := model.DateRgnOnDay(date)
	return es.repo.GetList(ctx, model.EventSearch{
		OwnerID:   &user.ID,
		DateRange: &dateRgn,
	})
}

func (es EventService) GetEventsOnWeek(ctx context.Context, user model.User, date time.Time) ([]model.Event, error) {
	dateRgn := model.DateRgnOnWeek(date)
	return es.repo.GetList(ctx, model.EventSearch{
		OwnerID:   &user.ID,
		DateRange: &dateRgn,
	})
}

func (es EventService) GetEventsOnMonth(ctx context.Context, user model.User, date time.Time) ([]model.Event, error) {
	dateRgn := model.DateRgnOnMonth(date)
	return es.repo.GetList(ctx, model.EventSearch{
		OwnerID:   &user.ID,
		DateRange: &dateRgn,
	})
}

func (es EventService) Delete(ctx context.Context, event model.Event) error {
	if err := es.repo.Delete(ctx, model.EventSearch{ID: &event.ID}); err != nil {
		return err
	}
	return nil
}

func (es EventService) GetByID(ctx context.Context, eventID uuid.UUID) (*model.Event, error) {
	return es.getOne(ctx, model.EventSearch{ID: &eventID})
}

func (es EventService) getOne(ctx context.Context, search model.EventSearch) (*model.Event, error) {
	events, err := es.repo.GetList(ctx, search)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, model.ErrEventNotFound
	}
	return &events[0], nil
}

func NewEventService(repo repository.Event, log logger.Logger) Event {
	return &EventService{
		repo: repo,
		log:  log,
	}
}
