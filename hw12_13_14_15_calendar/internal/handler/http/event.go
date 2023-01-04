package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	rs "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/rest/rqres"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

type Events struct {
	*Handler
}

func (e *Events) GetListOnDate(request *rs.Request) rs.Response {
	const actionName = "получение списка событий"
	// проверка правильности date и rangeType будет в сервисе
	rangeType, err := model.ParseRangeType(request.Param("rangeType"))
	if err != nil {
		e.handleError(actionName, fmt.Errorf("неверный rangeType %w", err))
	}
	date, err := time.Parse(time.RFC3339, request.URL.Query().Get("date"))
	if err != nil {
		e.handleError(actionName, fmt.Errorf("неверная дата %w", err))
	}
	events, err := e.services.Event.GetUserEventsOn(request.Context(), date, rangeType)
	if err != nil {
		err = fmt.Errorf("error events quering: %w", err)
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	return rs.Data(dto.FromEventSlice(events))
}

func (e *Events) GetByID(request *rs.Request) rs.Response {
	const actionName = "получение события по ID"
	ctx := request.Context()
	eventID, err := uuid.Parse(request.Param("eventID"))
	if err != nil {
		return e.handleError(actionName, fmt.Errorf("неверный eventID: %w", err))
	}
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		return e.handleError(actionName, err)
	}
	return rs.Data(dto.FromEventModel(*event))
}

func (e *Events) Create(request *rs.Request) rs.Response {
	const actionName = "добавление события"
	var input dto.EventCreate
	if request.ContentLength > 0 {
		defer func() {
			if err := request.Body.Close(); err != nil {
				e.logger.Error("добавление события - request.Body.Close(): %s", err.Error())
			}
		}()
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			return e.handleError(actionName, fmt.Errorf("ошибка парсинга входных данных: %w", err))
		}
	}
	inputCreate, vErrs := input.Model()
	if vErrs != nil {
		err := errx.InvalidNew("неверные данные", vErrs)
		return e.handleError(actionName, err)
	}
	event, err := e.services.Event.Add(request.Context(), inputCreate)
	if err != nil {
		err := fmt.Errorf("ошибка добавления события: %w", err)
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	e.logger.Info("событие добавлено: eventID=%s", event.ID.String())
	return rs.OK("событие добавлено", dto.FromEventModel(*event))
}

func (e *Events) Update(request *rs.Request) rs.Response {
	const actionName = "изменение события"
	var input dto.EventUpdate
	eventID, err := uuid.Parse(request.Param("eventID"))
	if err != nil {
		return e.handleError(actionName, fmt.Errorf("неверный eventID: %w", err))
	}
	ctx := request.Context()
	if request.ContentLength > 0 {
		defer func() {
			if err := request.Body.Close(); err != nil {
				e.logger.Error("изменение события - request.Body.Close(): %s", err.Error())
			}
		}()
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			return e.handleError(actionName, fmt.Errorf("ошибка парсинга входных данных: %w", err))
		}
	}
	inputUpdate, vErrs := input.Model()
	if vErrs != nil {
		err = errx.InvalidNew("неверные данные", vErrs)
		return e.handleError(actionName, err)
	}
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		return e.handleError(actionName, err)
	}
	err = e.services.Event.Update(request.Context(), *event, inputUpdate)
	if err != nil {
		return e.handleError(actionName, err)
	}
	e.logger.Info("событие изменено: eventID=%s", event.ID.String())
	return rs.OK("событие изменено", nil)
}

func (e *Events) Delete(request *rs.Request) rs.Response {
	const actionName = "удаление события"
	eventID, err := uuid.Parse(request.Param("eventID"))
	if err != nil {
		return e.handleError(actionName, fmt.Errorf("неверный eventID: %w", err))
	}
	ctx := request.Context()
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		return e.handleError(actionName, err)
	}
	err = e.services.Event.Delete(request.Context(), *event)
	if err != nil {
		return e.handleError(actionName, err)
	}
	e.logger.Info("событие удалено: eventID=%s", event.ID.String())
	return rs.OK("событие удалено", dto.FromEventModel(*event))
}
