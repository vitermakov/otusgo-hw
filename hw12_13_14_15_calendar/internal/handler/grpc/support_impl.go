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
		return nil, e.handleError(fmt.Errorf("ошибка запроса оповещений: %w", err))
	}
	return dto.FromNotificationSlice(notifiesList), nil
}

func (e SupportHandlerImpl) SetNotified(ctx context.Context, idReq *events.NotificationIDReq) (*emptypb.Empty, error) {
	eventID, err := dto.NotificationIDReqModel(idReq)
	if err != nil {
		return nil, e.handleError(fmt.Errorf("неверный идентификатор события: %w", err))
	}
	err = e.services.EventNotify.MarkEventNotified(ctx, eventID)
	if err != nil {
		return nil, e.handleError(fmt.Errorf("ошибка подтверждения оповещения события: %w", err))
	}
	e.logger.Info("событие изменено: eventID=%s", eventID.String())
	return &emptypb.Empty{}, nil
}

func (e SupportHandlerImpl) CleanupOldEvents(
	ctx context.Context, cleanupReq *events.CleanupReq,
) (*emptypb.Empty, error) {
	n, err := e.services.EventClean.CleanupOldEvents(ctx, cleanupReq.StoreTime.AsDuration())
	if err != nil {
		return nil, e.handleError(fmt.Errorf("ошибка удаления старых события: %w", err))
	}
	e.logger.Info("успешно удалено старых событий: %d", n)
	return &emptypb.Empty{}, nil
}

func (e SupportHandlerImpl) handleError(err error) error {
	e.logger.Error(err.Error())
	s := rqres.FromError(err)
	return status.Error(s.Code(), s.Message())
}
