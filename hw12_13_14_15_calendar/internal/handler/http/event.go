package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	rs "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest/rqres"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

type Events struct {
	*Handler
}

func (e *Events) ListOnDate(request *rs.Request) rs.Response {
	// проверка правильности rangeType будет в сервисе
	rangeType, _ := model.ParseRangeType(request.Param("rangeType"))
	date, _ := time.Parse(time.RFC3339, request.URL.Query().Get("date"))
	events, err := e.services.Event.GetUserEventsOn(request.Context(), date, rangeType)
	if err != nil {
		err = fmt.Errorf("error events quering: %w", err)
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	return rs.Data(dto.FromEventSlice(events))
}

func (e *Events) GetItem(request *rs.Request) rs.Response {
	ctx := request.Context()
	eventID, _ := uuid.Parse(request.Param("eventID"))
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	return rs.Data(dto.FromEventModel(*event))
}

func (e *Events) Create(request *rs.Request) rs.Response {
	var input dto.EventCreate
	if request.ContentLength > 0 {
		defer func() {
			if err := request.Body.Close(); err != nil {
				e.logger.Error("добавление события - request.Body.Close(): %s", err.Error())
			}
		}()
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			logErr := fmt.Errorf("добавление события - неверные входные данные: %w", err)
			e.logger.Error(logErr.Error())
			return rs.FromError(errx.LogicNew(logErr, 1000))
		}
	}
	event, err := e.services.Event.Add(request.Context(), input.Model())
	if err != nil {
		err := fmt.Errorf("ошибка добавления события: %w", err)
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	e.logger.Info("событие добавлено: eventID=%s", event.ID.String())
	return rs.OK("событие добавлено", dto.FromEventModel(*event))
}

func (e *Events) Update(request *rs.Request) rs.Response {
	var input dto.EventUpdate
	eventID, _ := uuid.Parse(request.Param("eventID"))
	ctx := request.Context()
	if request.ContentLength > 0 {
		defer func() {
			if err := request.Body.Close(); err != nil {
				e.logger.Error("изменение события - request.Body.Close(): %s", err.Error())
			}
		}()
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			logErr := fmt.Errorf("изменение события - неверные входные данные: %w", err)
			e.logger.Error(logErr.Error())
			return rs.FromError(errx.LogicNew(logErr, 1000))
		}
	}
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	err = e.services.Event.Update(request.Context(), *event, input.Model())
	if err != nil {
		err := fmt.Errorf("ошибка изменения события: %w", err)
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	e.logger.Info("событие изменено: eventID=%s", event.ID.String())
	return rs.OK("событие изменено", nil)
}

func (e *Events) Delete(request *rs.Request) rs.Response {
	eventID, _ := uuid.Parse(request.Param("event_id"))
	ctx := request.Context()
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	err = e.services.Event.Delete(request.Context(), *event)
	if err != nil {
		err := fmt.Errorf("ошибка удаления события: %w", err)
		e.logger.Error(err.Error())
		return rs.FromError(err)
	}
	e.logger.Info("событие удалено: eventID=%s", event.ID.String())
	return rs.OK("событие удалено", dto.FromEventModel(*event))
}
