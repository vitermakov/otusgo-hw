package scheduler

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
)

type Deps struct {
	Logger    logger.Logger
	API       *API
	Publisher queue.Producer
}

type API struct {
	Support events.SupportClient
}
