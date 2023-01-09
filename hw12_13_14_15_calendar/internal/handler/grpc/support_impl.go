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

func (e SupportHandlerImpl) GetNotifications(ctx context.Context, _ *emptypb.Empty) (*events.Notifies, error) {
	notifiesList, err := e.services.EventNotify.GetNotifications(ctx)
	if err != nil {
		err = fmt.Errorf("error events quering: %w", err)
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	return dto.FromNotificationSlice(notifiesList), nil
}

func (e SupportHandlerImpl) SetNotified(ctx context.Context, IDReq *events.NotificationIDReq) (*emptypb.Empty, error) {
	eventID := dto.NotificationIDReqModel(IDReq)
	err := e.services.EventNotify.MarkEventNotified(ctx, eventID)
	if err != nil {
		err := fmt.Errorf("ошибка подтверждения оповещения события: %w", err)
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	e.logger.Info("событие изменено: eventID=%s", eventID.String())
	return &emptypb.Empty{}, nil
}

func (e SupportHandlerImpl) CleanupOldEvents(ctx context.Context, cleanupReq *events.CleanupReq) (*emptypb.Empty, error) {
	n, err := e.services.EventClean.CleanupOldEvents(ctx, cleanupReq.StoreTime.AsDuration())
	if err != nil {
		err := fmt.Errorf("ошибка удаления старых события: %w", err)
		e.logger.Error(err.Error())
		s := rqres.FromError(err)
		return nil, status.Error(s.Code(), s.Message())
	}
	e.logger.Info("успешно удалено старых событий: %d", n)
	return &emptypb.Empty{}, nil
}
