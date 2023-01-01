package dto

import (
	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

/*
Тут оставим только преобразование из model.event.* и обратно
*/

type EventCreate events.CreateEvent
type EventUpdate events.UpdateEvent
type EventIDReq events.EventIDReq
type ListOnDateReq events.ListOnDateReq

// Model возвращает связанную модель если какие-то поля заполены
// в неверном формате, то не возвращаем ошибку, ошибка проверяется в сервисах.
func (ec *EventCreate) Model() model.EventCreate {
	// если ec == nil, возвращаем пустой model.EventCreate{}
	if ec == nil {
		return model.EventCreate{}
	}
	input := model.EventCreate{
		Title:       ec.Title,
		Duration:    ec.Duration.AsDuration(),
		Date:        ec.Date.AsTime(),
		Description: ec.Description,
	}
	if ec.Date != nil {
		input.Date = ec.Date.AsTime()
	}
	if ec.Duration != nil {
		input.Duration = ec.Duration.AsDuration()
	}
	if ec.Description != nil {
		val := *ec.Description
		input.Description = &val
	}
	if ec.NotifyTerm != nil {
		val := ec.NotifyTerm.AsDuration()
		input.NotifyTerm = &val
	}
	return input
}

func (eu *EventUpdate) Model() model.EventUpdate {
	// если ec == nil, возвращаем пустой model.EventCreate{}
	if eu == nil {
		return model.EventUpdate{}
	}
	input := model.EventUpdate{}
	if eu.Title != nil {
		val := *eu.Title
		input.Title = &val
	}
	if eu.Date != nil {
		// ошибки здесь не проверяем.
		val := eu.Date.AsTime()
		input.Date = &val
	}
	if eu.Duration != nil {
		val := eu.Duration.AsDuration()
		input.Duration = &val
	}
	if eu.Description != nil {
		val := *eu.Description
		input.Description = &val
	}
	if eu.NotifyTerm != nil {
		val := eu.NotifyTerm.AsDuration()
		input.NotifyTerm = &val
	}
	return input
}

func (eq *EventIDReq) Model() uuid.UUID {
	if eq == nil {
		return uuid.UUID{}
	}
	guid, _ := uuid.Parse(eq.ID)
	return guid
}

func (lr *ListOnDateReq) Model() (time.Time, model.RangeKind) {
	var (
		date      time.Time
		rangeType model.RangeKind
	)
	if lr == nil {
		return date, rangeType
	}
	if lr.Date != nil {
		date = lr.Date.AsTime()
	}
	rangeType = model.RangeKind(lr.RangeType.Number())

	return date, rangeType
}

func FromEventModel(item model.Event) *events.Event {
	return &events.Event{
		ID:          item.ID.String(),
		Title:       item.Title,
		Date:        timestamppb.New(item.Date),
		Duration:    durationpb.New(item.Duration),
		Description: item.Description,
		NotifyTerm:  durationpb.New(item.NotifyTerm),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
	}
}

func FromEventSlice(items []model.Event) *events.Events {
	result := &events.Events{
		List: nil,
	}
	if items == nil || len(items) == 0 {
		return result
	}
	result.List = make([]*events.Event, len(items))
	for i, item := range items {
		result.List[i] = FromEventModel(item)
	}
	return result
}
