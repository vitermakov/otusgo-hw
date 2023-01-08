package service

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type EventNotifyService struct {
	repo  repository.Event
	log   logger.Logger
	clock clock.Clock
}

func (e EventNotifyService) GetNotifications(ctx context.Context) ([]model.Notification, error) {
	return nil, nil
}

func (e EventNotifyService) MarkEventsNotified(ctx context.Context, eventID uuid.UUID) error {
	return nil
}

func NewEventNotifyService(repo repository.Event, log logger.Logger, clock clock.Clock) EventNotify {
	return &EventNotifyService{
		repo:  repo,
		log:   log,
		clock: clock,
	}
}
