package scheduler

import (
	"context"
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Actionable interface {
	DoAction(context.Context)
}

type Repeated struct {
	service Actionable
	period  time.Duration
	logger  logger.Logger
}

func NewRepeated(service Actionable, period time.Duration, logger logger.Logger) *Repeated {
	return &Repeated{service, period, logger}
}

func (r Repeated) Repeat(ctx context.Context) {
	go func() {
		t := time.NewTicker(r.period)
		defer func() {
			t.Stop()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				select {
				case <-ctx.Done():
					return
				default:
					r.service.DoAction(ctx)
				}
			}
		}
	}()
}
