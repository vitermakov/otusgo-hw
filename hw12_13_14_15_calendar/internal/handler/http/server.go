package http

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/rest"
)

func NewHandledServer(config config.Server, services *deps.Services, deps *deps.Deps) *rest.Server {
	server := rest.NewServer(servers.NewConfig(
		config.Host,
		config.Port,
		false,
	), services.Auth, deps.Logger)

	hs := NewHandlers(services, deps.Logger)

	server.GET("/events/list/{rangeType}", hs.Events.GetListOnDate)
	server.GET("/events/{eventID}", hs.Events.GetByID)
	server.POST("/events", hs.Events.Create)
	server.PUT("/events/{eventID}", hs.Events.Update)
	server.DELETE("/events/{eventID}", hs.Events.Delete)

	return server
}
