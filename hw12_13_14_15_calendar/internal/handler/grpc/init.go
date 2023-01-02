package grpc

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	grpcServ "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/grpc"
	"google.golang.org/grpc"
)

func InitHandlers(server *grpcServ.Server, services *deps.Services, dependencies *deps.Deps) {
	server.RegisterHandler(func(s *grpc.Server) {
		events.RegisterEventsServer(s, EventHandlerImpl{services: services, logger: dependencies.Logger})
	})
}
