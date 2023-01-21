package service

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

type EventNotifyService struct {
	repo  repository.Event
	log   logger.Logger
	clock clock.Clock
}

func (en EventNotifyService) GetNotifications(ctx context.Context) ([]model.Notification, error) {
	// events, err := en.repo.BlockEvents4Notify(ctx, en.clock.Now())
	now := en.clock.Now()
	events, err := en.repo.GetList(ctx, model.EventSearch{NeedNotifyTerm: &now})
	if err != nil {
		return nil, errx.FatalNew(err)
	}
	status := model.NotifyStatusBlocked
	_, err = en.repo.Update(ctx, model.EventUpdate{NotifyStatus: &status}, model.EventSearch{NeedNotifyTerm: &now})
	if err != nil {
		return nil, errx.FatalNew(err)
	}
	result := make([]model.Notification, len(events))
	for i, event := range events {
		user := model.NotifyUser{}
		if event.Owner != nil {
			user.Name = event.Owner.Name
			user.Email = event.Owner.Email
		}
		result[i] = model.Notification{
			EventID:       event.ID,
			EventTitle:    event.Title,
			EventDate:     event.Date,
			EventDuration: event.Duration,
			NotifyUser:    user,
		}
	}
	return result, nil
}

func (en EventNotifyService) MarkEventNotified(ctx context.Context, eventID uuid.UUID) error {
	nf := model.NotifyStatusNotified
	if _, err := en.repo.Update(ctx, model.EventUpdate{
		NotifyStatus: &nf,
	}, model.EventSearch{ID: &eventID}); err != nil {
		return errx.FatalNew(err)
	}
	return nil
}

func NewEventNotifyService(repo repository.Event, log logger.Logger, clock clock.Clock) EventNotify {
	return &EventNotifyService{
		repo:  repo,
		log:   log,
		clock: clock,
	}
}
