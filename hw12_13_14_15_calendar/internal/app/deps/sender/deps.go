package sender

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/mailer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
)

type Deps struct {
	Logger   logger.Logger
	API      *API
	Listener queue.Consumer
	Mailer   mailer.Mailer
}

type API struct {
	Support events.SupportClient
}
