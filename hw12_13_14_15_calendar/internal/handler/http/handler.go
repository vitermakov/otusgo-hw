package handler

import (
	"fmt"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/rest"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/rest/rqres"
)

type Handler struct {
	services *deps.Services
	logger   logger.Logger
}

func (h Handler) handleError(action string, err error) rqres.Response {
	logErr := fmt.Errorf("%s - %w", action, err)
	h.logger.Error(logErr.Error())
	return rqres.FromError(err)
}

type Handlers struct {
	Events *Events
}

func NewHandlers(services *deps.Services, logger logger.Logger) *Handlers {
	return &Handlers{
		Events: &Events{&Handler{services: services, logger: logger}},
	}
}

func (h *Handlers) InitRoutes(server *rest.Server) {
	server.GET("/events/list/{rangeType}", h.Events.GetListOnDate)
	server.GET("/events/{eventID}", h.Events.GetByID)
	server.POST("/events", h.Events.Create)
	server.PUT("/events/{eventID}", h.Events.Update)
	server.DELETE("/events/{eventID}", h.Events.Delete)
}
