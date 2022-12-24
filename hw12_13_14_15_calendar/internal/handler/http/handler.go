package handler

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest"
)

type Handler struct {
	services *deps.Services
	logger   logger.Logger
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
	server.GET("/events/list/on_day", h.Events.ListOnDay)
	server.GET("/events/list/on_week", h.Events.ListOnWeek)
	server.GET("/events/list/on_month", h.Events.ListOnMonth)
	server.GET("/events/:event_id", h.Events.GetItem)
	server.POST("/events", h.Events.Create)
	server.PUT("/events/:event_id", h.Events.Update)
	server.DELETE("/events/:event_id", h.Events.Delete)
}
