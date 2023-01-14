package scheduler

import (
	"context"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Notifier struct {
	supportAPI events.SupportClient
	authAPI    grpc.AuthFn
	publisher  queue.Producer
	logger     logger.Logger

	queueName string
}

func NewNotifier(
	api events.SupportClient, authAPI grpc.AuthFn, publisher queue.Producer, logger logger.Logger, queueName string,
) *Notifier {
	return &Notifier{
		supportAPI: api,
		authAPI:    authAPI,
		publisher:  publisher,
		logger:     logger,
		queueName:  queueName,
	}
}

func (ns Notifier) DoAction(ctx context.Context) {
	notificationsPb, err := ns.supportAPI.GetNotifications(ns.authAPI(ctx), &emptypb.Empty{})
	if err != nil {
		ns.logger.Error("notifier error getting notifications: %s", err.Error())
		return
	}
	notifications, err := dto.ToNotificationSlice(notificationsPb)
	if err != nil {
		ns.logger.Error("notifier error wrong data: %s", err.Error())
		return
	}
	if notifications == nil {
		ns.logger.Info("notifier: no new notifications")
		return
	}
	for _, note := range notifications {
		note := note
		message, err := queue.EncMessage(&note)
		if err != nil {
			ns.logger.Info("notifier error encoding notification: %s", err.Error())
			return
		}
		err = ns.publisher.Produce(message)
		if err != nil {
			ns.logger.Info("notifier error sending notification: %s", err.Error())
			return
		}
	}
	ns.logger.Info("notifier: %d notifications sent", len(notifications))
}
