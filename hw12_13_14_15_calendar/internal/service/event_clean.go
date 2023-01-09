package service

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
	"time"
)

type EventCleanService struct {
	repo  repository.Event
	log   logger.Logger
	clock clock.Clock
}

func (ec EventCleanService) CleanupOldEvents(ctx context.Context, timeLive time.Duration) (int64, error) {
	dateLess := ec.clock.Now().Add(timeLive * -1)
	n, err := ec.repo.Delete(ctx, model.EventSearch{DateLess: &dateLess})
	if err != nil {
		return 0, errx.FatalNew(err)
	}
	return n, nil
}

func NewEventCleanService(repo repository.Event, log logger.Logger, clock clock.Clock) EventClean {
	return &EventCleanService{
		repo:  repo,
		log:   log,
		clock: clock,
	}
}
