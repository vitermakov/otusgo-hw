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

// SupportHandlerImpl расширение генерированного GRPC сервера, внутренние запросы.
type SupportHandlerImpl struct {
	events.UnimplementedSupportServer
	services *deps.Services
	logger   logger.Logger
}

func (e SupportHandlerImpl) Update(ctx context.Context, updateEvent *events.UpdateEvent) (*emptypb.Empty, error) {
	eventID, input := dto.EventUpdateModel(updateEvent)
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	err = e.services.Event.Update(ctx, *event, input)
	if err != nil {
		err := fmt.Errorf("ошибка изменения события: %w", err)
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	e.logger.Info("событие изменено: eventID=%s", event.ID.String())
	return &emptypb.Empty{}, nil
}

func (e SupportHandlerImpl) Delete(ctx context.Context, idReq *events.EventIDReq) (*emptypb.Empty, error) {
	eventID := dto.EventIDReqModel(idReq)
	event, err := e.services.Event.GetByID(ctx, eventID)
	if err != nil {
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	err = e.services.Event.Delete(ctx, *event)
	if err != nil {
		err := fmt.Errorf("ошибка удаления события: %w", err)
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	e.logger.Info("событие удалено: eventID=%s", event.ID.String())
	return &emptypb.Empty{}, nil
}

func (e SupportHandlerImpl) GetListOnDate(ctx context.Context, lodReq *events.ListOnDateReq) (*events.Events, error) {
	date, rangeType := dto.ListOnDateReqModel(lodReq)
	evList, err := e.services.Event.GetUserEventsOn(ctx, date, rangeType)
	if err != nil {
		err = fmt.Errorf("error events quering: %w", err)
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	return dto.FromEventSlice(evList), nil
}
