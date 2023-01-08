package service

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type EventCleanService struct {
	repo  repository.Event
	log   logger.Logger
	clock clock.Clock
}

func (e EventCleanService) CleanupOldEvents(ctx context.Context) error {
	return nil
}

func NewEventCleanService(repo repository.Event, log logger.Logger, clock clock.Clock) EventClean {
	return &EventCleanService{
		repo:  repo,
		log:   log,
		clock: clock,
	}
}
