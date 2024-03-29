package grpc

import (
	"context"
	"fmt"

	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/grpc/rqres"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// EventHandlerImpl расширение генерированного GRPC сервера - для публичных запросов.
type EventHandlerImpl struct {
	events.UnimplementedEventsServer
	services *deps.Services
	logger   logger.Logger
}

func (e EventHandlerImpl) Create(ctx context.Context, createEvent *events.CreateEvent) (*events.Event, error) {
	event, err := e.services.EventCRUD.Add(ctx, dto.EventCreateModel(createEvent))
	if err != nil {
		return nil, e.handleError(fmt.Errorf("ошибка добавления события: %w", err))
	}
	e.logger.Info("событие добавлено: eventID=%s", event.ID.String())
	return dto.FromEventModel(*event), nil
}

func (e EventHandlerImpl) Update(ctx context.Context, updateEvent *events.UpdateEvent) (*emptypb.Empty, error) {
	eventID, input, err := dto.EventUpdateModel(updateEvent)
	if err != nil {
		return nil, e.handleError(fmt.Errorf("неверные данные обновления: %w", err))
	}
	event, err := e.services.EventCRUD.GetByID(ctx, eventID)
	if err != nil {
		return nil, e.handleError(err)
	}
	err = e.services.EventCRUD.Update(ctx, *event, input)
	if err != nil {
		return nil, e.handleError(fmt.Errorf("ошибка изменения события: %w", err))
	}
	e.logger.Info("событие изменено: eventID=%s", event.ID.String())
	return &emptypb.Empty{}, nil
}

func (e EventHandlerImpl) Delete(ctx context.Context, idReq *events.EventIDReq) (*emptypb.Empty, error) {
	eventID, err := dto.EventIDReqModel(idReq)
	if err != nil {
		return nil, e.handleError(fmt.Errorf("неверный идентификатор события: %w", err))
	}
	event, err := e.services.EventCRUD.GetByID(ctx, eventID)
	if err != nil {
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	err = e.services.EventCRUD.Delete(ctx, *event)
	if err != nil {
		err := fmt.Errorf("ошибка удаления события: %w", err)
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	e.logger.Info("событие удалено: eventID=%s", event.ID.String())
	return &emptypb.Empty{}, nil
}

func (e EventHandlerImpl) GetByID(ctx context.Context, idReq *events.EventIDReq) (*events.Event, error) {
	eventID, err := dto.EventIDReqModel(idReq)
	if err != nil {
		return nil, e.handleError(fmt.Errorf("неверный идентификатор события: %w", err))
	}
	event, err := e.services.EventCRUD.GetByID(ctx, eventID)
	if err != nil {
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	return dto.FromEventModel(*event), nil
}

func (e EventHandlerImpl) GetListOnDate(ctx context.Context, lodReq *events.ListOnDateReq) (*events.Events, error) {
	date, rangeType := dto.ListOnDateReqModel(lodReq)
	evList, err := e.services.EventCRUD.GetUserEventsOn(ctx, date, rangeType)
	if err != nil {
		return nil, e.handleError(fmt.Errorf("ошибка получения событий: %w", err))
	}
	return dto.FromEventSlice(evList), nil
}

func (e EventHandlerImpl) handleError(err error) error {
	e.logger.Error(err.Error())
	s := rqres.FromError(err)
	return status.Error(s.Code(), s.Message())
}
