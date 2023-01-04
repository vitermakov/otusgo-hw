package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

type EventService struct {
	repo repository.Event
	log  logger.Logger
	user User
}

func (es EventService) validateAdd(ctx context.Context, input model.EventCreate) error {
	if err := input.Validate(); err != nil {
		return err
	}
	user, err := es.user.GetByID(ctx, input.OwnerID)
	if err != nil {
		return err
	}
	if user == nil {
		return errx.LogicNew(model.ErrEventOwnerExists, model.ErrEventOwnerExistsCode)
	}
	events, err := es.repo.GetList(ctx, model.EventSearch{
		OwnerID: &input.OwnerID,
		DateRange: &model.DateRange{
			DateStart: input.Date,
			Duration:  input.Duration,
		},
		TacDuration: true,
	})
	if err != nil {
		return errx.FatalNew(err)
	}
	if len(events) > 0 {
		return errx.LogicNew(model.ErrEventDateBusy, model.ErrEventDateBusyCode)
	}
	return nil
}

func (es EventService) Add(ctx context.Context, input model.EventCreate) (*model.Event, error) {
	user, err := es.getAuthorizedUser(ctx, nil)
	if err != nil {
		return nil, err
	}
	input.OwnerID = user.ID
	if err = es.validateAdd(ctx, input); err != nil {
		errs := errx.ValidationErrors{}
		if errors.As(err, &errs) {
			return nil, errx.InvalidNew("неверные параметры", errs)
		}
		return nil, err
	}
	event, err := es.repo.Add(ctx, input)
	if err != nil {
		return nil, errx.FatalNew(err)
	}
	return event, nil
}

func (es EventService) validateUpdate(ctx context.Context, event model.Event, input model.EventUpdate) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if input.Date == nil && input.Duration == nil {
		return nil
	}
	dateRgn := model.DateRange{
		DateStart: event.Date,
		Duration:  event.Duration,
	}
	if input.Date != nil {
		dateRgn.DateStart = *input.Date
	}
	if input.Duration != nil {
		dateRgn.Duration = *input.Duration
	}
	search := model.EventSearch{
		OwnerID:     &event.Owner.ID,
		NotID:       &event.ID,
		DateRange:   &dateRgn,
		TacDuration: true,
	}
	events, err := es.repo.GetList(ctx, search)
	if err != nil {
		return errors.Wrap(errx.FatalNew(err), "ошибка проверки события")
	}
	if len(events) > 0 {
		return errx.LogicNew(model.ErrEventDateBusy, model.ErrEventDateBusyCode)
	}
	return nil
}

func (es EventService) Update(ctx context.Context, event model.Event, input model.EventUpdate) error {
	_, err := es.getAuthorizedUser(ctx, event.Owner)
	if err != nil {
		return err
	}
	if err = es.validateUpdate(ctx, event, input); err != nil {
		errs := errx.ValidationErrors{}
		if errors.As(err, &errs) {
			return errx.InvalidNew("неверные параметры", errs)
		}
		return err
	}
	if err = es.repo.Update(ctx, input, model.EventSearch{ID: &event.ID}); err != nil {
		return errx.FatalNew(err)
	}
	return nil
}

func (es EventService) GetUserEventsOn(
	ctx context.Context,
	date time.Time,
	kind model.RangeKind,
) ([]model.Event, error) {
	user, err := es.getAuthorizedUser(ctx, nil)
	if err != nil {
		return nil, err
	}
	dateRgn := model.DateRgnOn(kind, date)
	if !dateRgn.Valid() {
		return nil, errx.LogicNew(model.ErrCalendarDateRange, model.ErrCalendarDateRangeCode)
	}
	return es.GetEvents(ctx, model.EventSearch{
		OwnerID:   &user.ID,
		DateRange: &dateRgn,
	})
}

func (es EventService) GetEvents(ctx context.Context, search model.EventSearch) ([]model.Event, error) {
	events, err := es.repo.GetList(ctx, search)
	if err != nil {
		// неустранимая пользователем ошибка.
		return nil, errx.FatalNew(err)
	}
	return events, nil
}

func (es EventService) Delete(ctx context.Context, event model.Event) error {
	_, err := es.getAuthorizedUser(ctx, event.Owner)
	if err != nil {
		return err
	}
	if err = es.repo.Delete(ctx, model.EventSearch{ID: &event.ID}); err != nil {
		// неустранимая пользователем ошибка.
		return errx.FatalNew(err)
	}
	return nil
}

func (es EventService) GetByID(ctx context.Context, eventID uuid.UUID) (*model.Event, error) {
	event, err := es.getOne(ctx, model.EventSearch{ID: &eventID})
	if err == nil {
		return event, nil
	}
	// если ошибка - NotFound, добавим параметр eventId.
	nfErr := errx.NotFound{}
	if errors.As(err, &nfErr) {
		nfErr.Params = map[string]uuid.UUID{
			"eventId": eventID,
		}
		return nil, nfErr
	}
	return nil, err
}

func (es EventService) getOne(ctx context.Context, search model.EventSearch) (*model.Event, error) {
	events, err := es.GetEvents(ctx, search)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, errx.NotFoundNew(model.ErrEventNotFound, nil)
	}
	return &events[0], nil
}

// getAuthorizedUser получить текущего пользователя.
func (es EventService) getAuthorizedUser(ctx context.Context, checkUser *model.User) (*model.User, error) {
	user, err := es.user.GetCurrent(ctx)
	if err != nil {
		// пользователь не авторизован.
		nfErr := errx.NotFound{}
		if errors.As(err, &nfErr) {
			return nil, errx.LogicNew(model.ErrCalendarAccess, model.ErrCalendarAccessCode)
		}
		return nil, errx.FatalNew(err)
	}
	if checkUser != nil && strings.Compare(user.ID.String(), checkUser.ID.String()) != 0 {
		return nil, errx.LogicNew(model.ErrCalendarAccess, model.ErrCalendarAccessCode)
	}
	return user, nil
}

func NewEventService(repo repository.Event, log logger.Logger, user User) Event {
	return &EventService{
		repo: repo,
		log:  log,
		user: user,
	}
}
