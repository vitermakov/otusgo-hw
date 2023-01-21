package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/client/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
)

type MainSuiteTest struct {
	suite.Suite
	client   calendar.Client
	eventIds []string
}

const (
	ValidUserEmail = "ivan@otus.ru"
	CalendarAPIURL = "http://localhost:8080"
)

func (ms *MainSuiteTest) SetupTest() {
	ms.client = calendar.NewClient(CalendarAPIURL, ValidUserEmail)
	ms.eventIds = make([]string, 0, 20)
}

func (ms *MainSuiteTest) TearDownTest() {
	ctx := context.Background()
	for _, id := range ms.eventIds {
		err := ms.client.Delete(ctx, id)
		ms.Suite.Require().NoError(err)
	}
}

func (ms *MainSuiteTest) TestCreate() {
	dateInFuture := time.Now().Add(time.Hour * 24 * 30)
	desc := "№348239"
	nfTermOk := "3h"
	nfTermErr := "6q"
	testCases := []struct {
		name          string
		inputCreate   dto.EventCreate
		checkResponse func(*dto.Event, error)
	}{
		{
			name: "event add ok",
			inputCreate: dto.EventCreate{
				Title:       "Встреча в Zoom",
				Date:        dateInFuture.Format(time.RFC3339),
				Duration:    "45m",
				Description: &desc,
				NotifyTerm:  &nfTermOk,
			},
			checkResponse: func(event *dto.Event, rErr error) {
				ms.Suite.Require().NoError(rErr)
				_, err := uuid.Parse(event.ID)
				ms.Suite.Require().NoError(err)
				d, _ := time.ParseDuration("45m")
				ms.Suite.Require().Equal(d.String(), event.Duration)
				nt, _ := time.ParseDuration("3h")
				ms.Suite.Require().Equal(nt.String(), event.NotifyTerm)
			},
		}, {
			name:        "event empty",
			inputCreate: dto.EventCreate{},
			checkResponse: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				invErr := errx.Invalid{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
				ms.Suite.Require().Len(invErr.Errors(), 2)
			},
		}, {
			name: "event wrong data",
			inputCreate: dto.EventCreate{
				Title:      "Test event",
				Date:       "----2023-02-19T20:00:00.417Z-----",
				Duration:   "45q",
				NotifyTerm: &nfTermErr,
			},
			checkResponse: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				invErr := errx.Invalid{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
				ms.Suite.Require().Len(invErr.Errors(), 3)
			},
		}, {
			name: "duplicate entry",
			inputCreate: dto.EventCreate{
				Title:    "Встреча в Zoom",
				Date:     dateInFuture.Add(time.Minute * 10).Format(time.RFC3339),
				Duration: "50m",
			},
			checkResponse: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
				ms.Suite.Require().Equal(model.ErrEventDateBusyCode, logErr.Code())
			},
		},
	}
	for _, tc := range testCases {
		ms.Suite.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			event, err := ms.client.Create(ctx, tc.inputCreate)
			if event != nil {
				ms.eventIds = append(ms.eventIds, event.ID)
			}
			tc.checkResponse(event, err)
		})
	}
}

func (ms *MainSuiteTest) TestGetByID() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dateInFuture := time.Now().Add(time.Hour * 24 * 30)
	desc := "№348239"
	nfTermOk := "3h"
	newEvent, err := ms.client.Create(ctx, dto.EventCreate{
		Title:       "Test event",
		Date:        dateInFuture.Format(time.RFC3339),
		Duration:    "45m",
		Description: &desc,
		NotifyTerm:  &nfTermOk,
	})
	ms.Suite.Require().NotNil(newEvent)
	ms.eventIds = append(ms.eventIds, newEvent.ID)
	ms.Suite.Require().NoError(err)

	testCases := []struct {
		name          string
		ID            string
		checkResponse func(event *dto.Event, err error)
	}{
		{
			name: "exists event",
			ID:   newEvent.ID,
			checkResponse: func(event *dto.Event, rErr error) {
				ms.Suite.Require().NotNil(event)
				ms.Suite.Require().NoError(rErr)
				ms.Suite.Require().Equal(newEvent.ID, event.ID)
			},
		}, {
			name: "wrong event ID",
			ID:   "ascascas",
			checkResponse: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
			},
		}, {
			name: "not exists event",
			ID:   uuid.New().String(),
			checkResponse: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				invErr := errx.NotFound{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		ms.Suite.Run(tc.name, func() {
			tc.checkResponse(ms.client.GetByID(ctx, tc.ID))
		})
	}
}

func (ms *MainSuiteTest) TestUpdate() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dateInFuture := time.Now().Add(time.Hour * 24 * 30)
	insertInput := []dto.EventCreate{
		{
			Title:    "Test event 1",
			Date:     dateInFuture.Format(time.RFC3339),
			Duration: "60m",
		}, {
			Title:    "Test event 2",
			Date:     dateInFuture.Add(time.Hour * 2).Format(time.RFC3339),
			Duration: "60m",
		},
	}
	events := make([]dto.Event, 0, 2)
	for _, input := range insertInput {
		newEvent, err := ms.client.Create(ctx, input)
		ms.Suite.Require().NotNil(newEvent)
		ms.eventIds = append(ms.eventIds, newEvent.ID)
		ms.Suite.Require().NoError(err)
		events = append(events, *newEvent)
	}
	titleOk := "Test event 2 updated"
	dateOk := dateInFuture.Add(time.Hour * 4).Format(time.RFC3339)
	durOk := "90m"
	desc := "new description"
	nfTermOk := "3h"

	dateDup := dateInFuture.Add(time.Minute * 30 * -1).Format(time.RFC3339)
	durDup := "60m"

	titleErr := ""
	dateErr := "-34t---"
	durErr := "45o"
	nfTermErr := "6o"

	testCases := []struct {
		name          string
		inputUpdate   dto.EventUpdate
		ID            string
		checkResponse func(error)
	}{
		{
			name: "event update ok",
			inputUpdate: dto.EventUpdate{
				Title:       &titleOk,
				Date:        &dateOk,
				Duration:    &durOk,
				Description: &desc,
				NotifyTerm:  &nfTermOk,
			},
			ID: events[1].ID,
			checkResponse: func(rErr error) {
				ms.Suite.Require().NoError(rErr)
				event, err := ms.client.GetByID(ctx, events[1].ID)
				ms.Suite.Require().NotNil(event)
				ms.Suite.Require().NoError(err)

				// убедиться что событие обновилось
				ms.Suite.Require().Equal(events[1].ID, event.ID)
				ms.Suite.Require().Equal(titleOk, event.Title)
				ms.Suite.Require().Equal(dateOk, event.Date.Format(time.RFC3339))
				dv, _ := time.ParseDuration(durOk)
				ms.Suite.Require().Equal(dv.String(), event.Duration)
				ms.Suite.Require().Equal(desc, event.Description)
				dv, _ = time.ParseDuration(nfTermOk)
				ms.Suite.Require().Equal(dv.String(), event.NotifyTerm)
			},
		}, {
			name: "event move to occupied date",
			inputUpdate: dto.EventUpdate{
				Date:     &dateDup,
				Duration: &durDup,
			},
			ID: events[1].ID,
			checkResponse: func(rErr error) {
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
				ms.Suite.Require().Equal(model.ErrEventDateBusyCode, logErr.Code())
			},
		}, {
			name: "event wrong data",
			inputUpdate: dto.EventUpdate{
				Title:      &titleErr,
				Date:       &dateErr,
				Duration:   &durErr,
				NotifyTerm: &nfTermErr,
			},
			ID: events[0].ID,
			checkResponse: func(rErr error) {
				invErr := errx.Invalid{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
				ms.Suite.Require().Len(invErr.Errors(), 3)
			},
		},
	}

	for _, tc := range testCases {
		ms.Suite.Run(tc.name, func() {
			tc.checkResponse(ms.client.Update(ctx, tc.ID, tc.inputUpdate))
		})
	}
}

func (ms *MainSuiteTest) TestDelete() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dateInFuture := time.Now().Add(time.Hour * 24 * 30)
	newEvent, err := ms.client.Create(ctx, dto.EventCreate{
		Title:    "Test event",
		Date:     dateInFuture.Format(time.RFC3339),
		Duration: "45m",
	})
	ms.Suite.Require().NotNil(newEvent)
	ms.Suite.Require().NoError(err)

	ms.Suite.Run("removing exists event", func() {
		err = ms.client.Delete(ctx, newEvent.ID)
		ms.Suite.Require().NoError(err)

		_, err = ms.client.GetByID(ctx, newEvent.ID)
		rErr := errx.NotFound{}
		ms.Suite.Require().ErrorAs(err, &rErr)
	})

	ms.Suite.Run("removing not exists event", func() {
		err = ms.client.Delete(ctx, newEvent.ID)
		rErr := errx.NotFound{}
		ms.Suite.Require().ErrorAs(err, &rErr)
	})
}

func (ms *MainSuiteTest) TestList() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	insertInput := []dto.EventCreate{
		{
			Title:    "Test event 1",
			Date:     "2023-02-19T20:00:00.417Z",
			Duration: "45m",
		}, {
			Title:    "Test event 2",
			Date:     "2023-02-19T09:00:00.417Z",
			Duration: "90m",
		}, {
			Title:    "Test event 3",
			Date:     "2023-02-13T15:00:00.417Z",
			Duration: "60m",
		}, {
			Title:    "Test event 4",
			Date:     "2023-03-10T10:00:00.417Z",
			Duration: "45m",
		}, {
			Title:    "Test event 5",
			Date:     "2023-03-22T20:00:00.417Z",
			Duration: "45m",
		},
	}
	events := make([]*dto.Event, len(insertInput))
	for i, input := range insertInput {
		event, err := ms.client.Create(ctx, input)
		ms.Suite.Require().NotNil(event)
		ms.Suite.Require().NoError(err)
		ms.eventIds = append(ms.eventIds, event.ID)
		events[i] = event
	}
	testCases := []struct {
		name          string
		rangeType     string
		date          string
		checkResponse func(error)
		expectedIDs   []string
	}{
		{
			name:      "day (whole 2023-02-19) events",
			rangeType: "day",
			date:      "2023-02-19T22:00:00.417Z",
			expectedIDs: []string{
				events[0].ID,
				events[1].ID,
			},
			checkResponse: func(rErr error) {
				ms.Suite.Require().NoError(rErr)
			},
		}, {
			name:      "month (whole march) events",
			rangeType: "month",
			date:      "2023-03-19T22:00:00.417Z",
			expectedIDs: []string{
				events[3].ID,
				events[4].ID,
			},
			checkResponse: func(rErr error) {
				ms.Suite.Require().NoError(rErr)
			},
		}, {
			name:      "week feb 13-19",
			rangeType: "week",
			date:      "2023-02-15T22:00:00.417Z",
			expectedIDs: []string{
				events[0].ID,
				events[1].ID,
				events[2].ID,
			},
			checkResponse: func(rErr error) {
				ms.Suite.Require().NoError(rErr)
			},
		}, {
			name:        "empty list",
			rangeType:   "week",
			date:        "2023-02-25T22:00:00.417Z",
			expectedIDs: []string{},
			checkResponse: func(rErr error) {
				ms.Suite.Require().NoError(rErr)
			},
		}, {
			name:        "wrong range type",
			rangeType:   "hour",
			date:        "2023-02-25T22:00:00.417Z",
			expectedIDs: []string{},
			checkResponse: func(rErr error) {
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
				ms.Suite.Require().Equal(model.ErrCalendarDateRangeCode, logErr.Code())
			},
		}, {
			name:        "wrong date",
			rangeType:   "day",
			date:        "-------2023-02-25T22:00:00.417Z-----",
			expectedIDs: []string{},
			checkResponse: func(rErr error) {
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
				ms.Suite.Require().Equal(model.ErrCalendarDateRangeCode, logErr.Code())
			},
		},
	}
	for _, tc := range testCases {
		ms.Suite.Run(tc.name, func() {
			date, _ := time.Parse(time.RFC3339, tc.date)
			events, err := ms.client.GetListOnDate(ctx, tc.rangeType, date)
			actualIDs := make([]string, len(events))
			for i, event := range events {
				actualIDs[i] = event.ID
			}
			ms.Suite.Require().Equal(tc.expectedIDs, actualIDs)

			tc.checkResponse(err)
		})
	}
}

func (ms *MainSuiteTest) TestNotify() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	dateNowDays := time.Now().Add(time.Minute * 20)
	dateIndaFuture := time.Now().Add(time.Hour * 2)
	notifyTerm := "30m"

	insertInput := []dto.EventCreate{
		{
			Title:      "Event must to be notified",
			Date:       dateNowDays.Format(time.RFC3339),
			Duration:   "45m",
			NotifyTerm: &notifyTerm,
		}, {
			Title:      "Event must not to be notified",
			Date:       dateIndaFuture.Format(time.RFC3339),
			Duration:   "45m",
			NotifyTerm: &notifyTerm,
		},
	}
	events := make([]*dto.Event, len(insertInput))
	for i, input := range insertInput {
		event, err := ms.client.Create(ctx, input)
		ms.Suite.Require().NotNil(event)
		ms.Suite.Require().NoError(err)
		ms.eventIds = append(ms.eventIds, event.ID)
		events[i] = event
	}

	testCases := []struct {
		name           string
		ID             string
		notifiedStatus model.NotifyStatus
	}{
		{
			name:           "event must to be notified",
			ID:             events[0].ID,
			notifiedStatus: model.NotifyStatusNotified,
		}, {
			name:           "Event must not to be notified",
			ID:             events[1].ID,
			notifiedStatus: model.NotifyStatusNone,
		},
	}

	// waiting for notifier trigger
	select {
	case <-ctx.Done():
		ms.Suite.Require().True(false, "testing adding or timer is too long")
	case <-time.After(time.Second * 7):
	}

	for _, tc := range testCases {
		ms.Suite.Run(tc.name, func() {
			event, err := ms.client.GetByID(ctx, tc.ID)
			ms.Suite.Require().NoError(err)
			ms.Suite.Require().NotNil(event)

			notifyStatus, err := model.ParseNotifyStatus(event.NotifyStatus)
			ms.Suite.Require().NoError(err)

			ms.Suite.Require().Equal(notifyStatus, tc.notifiedStatus)
		})
	}
}

func TestCalendar(t *testing.T) {
	suite.Run(t, new(MainSuiteTest))
}
