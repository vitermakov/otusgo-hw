package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
)

type Notifier struct {
	supportAPI events.SupportClient
	publisher  queue.Producer
	logger     logger.Logger

	queueName string
}

func NewNotifier(
	api events.SupportClient, publisher queue.Producer, logger logger.Logger, queueName string,
) *Notifier {
	return &Notifier{
		supportAPI: api,
		publisher:  publisher,
		logger:     logger,
		queueName:  queueName,
	}
}

func (ns Notifier) DoAction(ctx context.Context) {
	/*
		notificationsPb, err := ns.supportAPI.GetNotifications(ctx, &emptypb.Empty{})
		if err != nil {
			ns.Logger.Error("notifier error getting notifications: %s", err.Error())
			return
		}
		notifications := dto.ToNotificationSlice(notificationsPb)
		if notifications == nil {
			ns.Logger.Info("notifier: no new notifications")
			return
		}
	*/
	notifications := []model.Notification{
		{
			EventID:       uuid.New(),
			EventTitle:    fmt.Sprintf("Event Title %d", time.Now().UnixMilli()%1000),
			EventDate:     time.Now(),
			EventDuration: time.Minute * 45,
			NotifyUser: model.NotifyUser{
				Name:  "Andrew",
				Email: "vit_ermakov@mail.ru",
			},
		},
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
