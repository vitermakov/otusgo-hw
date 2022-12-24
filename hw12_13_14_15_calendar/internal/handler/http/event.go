package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
	rs "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest/rqres"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
	"time"
)

type Events struct {
	*Handler
}

func (e *Events) ListOnDay(request *rs.Request) rs.Response {
	events, err := e.services.Event.GetEventsOnDay(request.Context(), time.Now())
	if err != nil {
		msg := fmt.Errorf("error events quering: %w", err)
		e.logger.Error(msg.Error())
		return rs.FromError(errx.LogicNew(msg, 1000))
	}
	return rs.Data(dto.FromEventSlice(events))
}

func (e *Events) ListOnWeek(request *rs.Request) rs.Response {
	events, err := e.services.Event.GetEventsOnWeek(request.Context(), time.Now())
	if err != nil {
		msg := fmt.Errorf("error events quering: %w", err)
		e.logger.Error(msg.Error())
		return rs.FromError(errx.LogicNew(msg, 1000))
	}
	return rs.Data(dto.FromEventSlice(events))
}

func (e *Events) ListOnMonth(request *rs.Request) rs.Response {
	events, err := e.services.Event.GetEventsOnMonth(request.Context(), time.Now())
	if err != nil {
		msg := fmt.Errorf("error events quering: %w", err)
		e.logger.Error(msg.Error())
		return rs.FromError(errx.LogicNew(msg, 1000))
	}
	return rs.Data(dto.FromEventSlice(events))
}

func (e *Events) GetItem(request *rs.Request) rs.Response {
	ctx := request.Context()
	eventId, _ := uuid.Parse(request.Param("event_id"))
	event, err := e.services.Event.GetByID(ctx, eventId)
	if err != nil {
		msg := fmt.Errorf("error event id='%s': %w", eventId.String(), err)
		e.logger.Error(msg.Error())
		return rs.FromError(err)
	}
	return rs.Data(dto.FromEventModel(*event))
}
func (e *Events) Create(request *rs.Request) rs.Response {
	var input dto.EventCreate
	if request.ContentLength > 0 {
		defer func() {
			// тут не знаю, это тоже надо логгировать?
			_ = request.Body.Close()
		}()
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			logErr := fmt.Errorf("error parsing input in event create: %w", err)
			e.logger.Error(logErr.Error())
			return rs.FromError(errx.LogicNew(logErr, 1000))
		}
	}
	event, err := e.services.Event.Add(request.Context(), input.Model())
	if err != nil {
		msg := fmt.Errorf("error parsing input in event create: %w", err)
		e.logger.Error(msg.Error())
		return rs.FromError(errx.LogicNew(msg, 1000))
	}
	return rs.OK("event created", dto.FromEventModel(*event))
}
func (e *Events) Update(request *rs.Request) rs.Response {
	var input dto.EventUpdate
	eventId, _ := uuid.Parse(request.Param("event_id"))
	ctx := request.Context()
	if request.ContentLength > 0 {
		defer func() {
			// тут не знаю, это тоже надо логгировать?
			_ = request.Body.Close()
		}()
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			logErr := fmt.Errorf("error parsing input in event update: %w", err)
			e.logger.Error(logErr.Error())
			return rs.FromError(errx.LogicNew(logErr, 1000))
		}
	}
	event, err := e.services.Event.GetByID(ctx, eventId)
	if err != nil {
		msg := fmt.Errorf("error event id='%s': %w", eventId.String(), err)
		e.logger.Error(msg.Error())
		return rs.FromError(err)
	}
	err = e.services.Event.Update(request.Context(), *event, input.Model())
	if err != nil {
		msg := fmt.Errorf("error parsing input in event update: %s", err.Error())
		e.logger.Error(msg.Error())
		return rs.FromError(errx.LogicNew(msg, 1000))
	}
	return rs.OK("event updated", dto.FromEventModel(*event))
}
func (e *Events) Delete(request *rs.Request) rs.Response {
	eventId, _ := uuid.Parse(request.Param("event_id"))
	ctx := request.Context()

	event, err := e.services.Event.GetByID(ctx, eventId)
	if err != nil {
		msg := fmt.Errorf("error event id='%s': %w", eventId.String(), err)
		e.logger.Error(msg.Error())
		return rs.FromError(err)
	}

	err = e.services.Event.Delete(request.Context(), *event)
	if err != nil {
		msg := fmt.Errorf("error parsing input in event update: %s", err.Error())
		e.logger.Error(msg.Error())
		return rs.FromError(errx.LogicNew(msg, 1000))
	}
	return rs.OK("event deleted", dto.FromEventModel(*event))
}
