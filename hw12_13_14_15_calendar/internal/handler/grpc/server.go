package grpc

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers"
	grpcServ "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/grpc"
	"google.golang.org/grpc"
)

func NewHandledServer(config config.Server, services *deps.Services, deps *deps.Deps) *grpcServ.Server {
	server := grpcServ.NewServer(servers.NewConfig(
		config.Host,
		config.Port,
		false,
	), services.Auth, deps.Logger)

	server.RegisterHandler(func(s *grpc.Server) {
		events.RegisterEventsServer(s, EventHandlerImpl{services: services, logger: deps.Logger})
		events.RegisterSupportServer(s, SupportHandlerImpl{services: services, logger: deps.Logger})
	})

	return server
}
