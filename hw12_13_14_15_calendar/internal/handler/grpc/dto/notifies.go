package dto

import (
	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NotificationIDReqModel(idReq *events.NotificationIDReq) uuid.UUID {
	if idReq == nil {
		return uuid.UUID{}
	}
	guid, _ := uuid.Parse(idReq.ID)
	return guid
}

func FromNotificationModel(item model.Notification) *events.Notification {
	return &events.Notification{
		ID:        item.EventID.String(),
		Title:     item.EventTitle,
		Date:      timestamppb.New(item.EventDate),
		Duration:  durationpb.New(item.EventDuration),
		UserName:  item.NotifyUser.Name,
		UserEmail: item.NotifyUser.Email,
	}
}

func FromNotificationSlice(items []model.Notification) *events.Notifies {
	result := &events.Notifies{
		List: nil,
	}
	if len(items) == 0 {
		return result
	}
	result.List = make([]*events.Notification, len(items))
	for i, item := range items {
		result.List[i] = FromNotificationModel(item)
	}
	return result
}

func ToNotificationModel(item *events.Notification) model.Notification {
	if item == nil {
		return model.Notification{}
	}
	eventID, _ := uuid.Parse(item.ID)
	return model.Notification{
		EventID:       eventID,
		EventTitle:    item.Title,
		EventDate:     item.Date.AsTime(),
		EventDuration: item.Duration.AsDuration(),
		NotifyUser: model.NotifyUser{
			Name:  item.UserName,
			Email: item.UserEmail,
		},
	}
}

func ToNotificationSlice(items *events.Notifies) []model.Notification {
	if items == nil {
		return nil
	}
	result := make([]model.Notification, len(items.List))
	for i, item := range items.List {
		result[i] = ToNotificationModel(item)
	}
	return result
}
