package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
Тут оставим только преобразование из model.event.* и обратно
*/

// EventCreateModel возвращает связанную модель если какие-то поля заполены
// в неверном формате, то не возвращаем ошибку, ошибка проверяется в сервисах.
func EventCreateModel(createEvent *events.CreateEvent) model.EventCreate {
	// если ec == nil, возвращаем пустой model.EventCreate{}
	if createEvent == nil {
		return model.EventCreate{}
	}
	input := model.EventCreate{
		Title:       createEvent.Title,
		Duration:    createEvent.Duration.AsDuration(),
		Date:        createEvent.Date.AsTime(),
		Description: createEvent.Description,
	}
	if createEvent.Date != nil {
		input.Date = createEvent.Date.AsTime()
	}
	if createEvent.Duration != nil {
		input.Duration = createEvent.Duration.AsDuration()
	}
	if createEvent.Description != nil {
		val := *createEvent.Description
		input.Description = &val
	}
	if createEvent.NotifyTerm != nil {
		val := createEvent.NotifyTerm.AsDuration()
		input.NotifyTerm = &val
	}
	return input
}

func EventUpdateModel(updateEvent *events.UpdateEvent) (uuid.UUID, model.EventUpdate, error) {
	// если ec == nil, возвращаем пустой model.EventUpdate{}
	if updateEvent == nil {
		return uuid.UUID{}, model.EventUpdate{}, errors.New("empty query")
	}
	input := model.EventUpdate{}
	if updateEvent.Title != nil {
		val := *updateEvent.Title
		input.Title = &val
	}
	if updateEvent.Date != nil {
		// ошибки здесь не проверяем.
		val := updateEvent.Date.AsTime()
		input.Date = &val
	}
	if updateEvent.Duration != nil {
		val := updateEvent.Duration.AsDuration()
		input.Duration = &val
	}
	if updateEvent.Description != nil {
		val := *updateEvent.Description
		input.Description = &val
	}
	if updateEvent.NotifyTerm != nil {
		val := updateEvent.NotifyTerm.AsDuration()
		input.NotifyTerm = &val
	}
	guid, err := uuid.Parse(updateEvent.ID)
	if err != nil {
		return uuid.UUID{}, model.EventUpdate{}, err
	}
	return guid, input, nil
}

func EventIDReqModel(idReq *events.EventIDReq) (uuid.UUID, error) {
	if idReq == nil {
		return uuid.UUID{}, errors.New("empty eventIDReq")
	}
	return uuid.Parse(idReq.ID)
}

func ListOnDateReqModel(req *events.ListOnDateReq) (time.Time, model.RangeKind) {
	var (
		date      time.Time
		rangeType model.RangeKind
	)
	if req == nil {
		return date, rangeType
	}
	if req.Date != nil {
		date = req.Date.AsTime()
	}
	rangeType = model.RangeKind(req.RangeType.Number())

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
	if len(items) == 0 {
		return result
	}
	result.List = make([]*events.Event, len(items))
	for i, item := range items {
		result.List[i] = FromEventModel(item)
	}
	return result
}
