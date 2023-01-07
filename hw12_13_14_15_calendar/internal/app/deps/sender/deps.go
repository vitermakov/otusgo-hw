package scheduler

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Deps struct {
	Logger logger.Logger
	APIs   *APIs
}

type Services struct {
}

type APIs struct {
}

func NewServices(deps *Deps) *Services {
	// api := deps.APIs

	return &Services{}
}
