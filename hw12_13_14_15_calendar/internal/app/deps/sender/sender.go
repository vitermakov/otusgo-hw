package sender

import (
	"context"
	"fmt"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/mailer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
)

type Sender struct {
	supportAPI events.SupportClient
	authAPI    grpc.AuthFn
	listener   queue.Consumer
	logger     logger.Logger
	mailer     mailer.Mailer

	queueName   string
	defaultFrom string
}

func NewSender(
	api events.SupportClient, authAPI grpc.AuthFn, consumer queue.Consumer, l logger.Logger, ml mailer.Mailer,
	qn, from string,
) *Sender {
	return &Sender{
		supportAPI: api, authAPI: authAPI, listener: consumer, logger: l, mailer: ml,
		queueName: qn, defaultFrom: from,
	}
}

func (s Sender) Run(ctx context.Context) error {
	msgChan, err := s.listener.Consume(ctx, s.queueName)
	if err != nil {
		return fmt.Errorf("error initializing consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-msgChan:
				var note model.Notification
				err := message.Decode(&note)
				if err != nil {
					s.logger.Error("sender can't parse notification: %s", err.Error())
					break
				}
				err = s.sendMailAndConfirm(ctx, note)
				if err != nil {
					s.logger.Error("sender can't sending notification: %s", err.Error())
					break
				}
				s.logger.Info(
					"notification event %s on %s sent to %s",
					note.EventID.String(), note.EventDate.String(), note.NotifyUser.Email,
				)
			}
		}
	}()
	return nil
}

func (s Sender) sendMailAndConfirm(ctx context.Context, note model.Notification) error {
	dateFmt := "02.01.2006 15:04"
	sendData := map[string]string{
		"UserName":       note.NotifyUser.Name,
		"UserEmail":      note.NotifyUser.Email,
		"EventTitle":     note.EventTitle,
		"EventID":        note.EventID.String(),
		"SenderEmail":    s.defaultFrom,
		"EventDateStart": note.EventDate.Format(dateFmt),
		"EventDateEnd":   note.EventDate.Add(note.EventDuration).Format(dateFmt),
	}
	err := s.mailer.SendMail("events/notify", mailer.Mail{
		Sender:  s.defaultFrom,
		To:      []string{note.NotifyUser.Email},
		Subject: fmt.Sprintf("%s: начнется %s", note.EventTitle, note.EventDate.Format(dateFmt)),
		Data:    sendData,
	})
	if err != nil {
		return err
	}
	_, err = s.supportAPI.SetNotified(s.authAPI(ctx), &events.NotificationIDReq{ID: note.EventID.String()})
	return err
}
