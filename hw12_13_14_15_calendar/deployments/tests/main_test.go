package main

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/client/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
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
		name              string
		inputCreate       dto.EventCreate
		checkResponseBody func(*dto.Event, error)
	}{
		{
			name: "event add ok",
			inputCreate: dto.EventCreate{
				Title:       "Встреча в Zoom",
				Date:        dateInFuture.String(),
				Duration:    "45m",
				Description: &desc,
				NotifyTerm:  &nfTermOk,
			},
			checkResponseBody: func(event *dto.Event, rErr error) {
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
			checkResponseBody: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				invErr := errx.Invalid{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
				ms.Suite.Require().Len(invErr.Errors, 2)
			},
		}, {
			name: "event wrong data",
			inputCreate: dto.EventCreate{
				Title:      "Test event",
				Date:       "----2023-02-19T20:00:00.417Z-----",
				Duration:   "45q",
				NotifyTerm: &nfTermErr,
			},
			checkResponseBody: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				invErr := errx.Invalid{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
				ms.Suite.Require().Len(invErr.Errors, 3)
			},
		}, {
			name: "duplicate entry",
			inputCreate: dto.EventCreate{
				Title:    "Встреча в Zoom",
				Date:     dateInFuture.Add(time.Minute * 10).String(),
				Duration: "50m",
			},
			checkResponseBody: func(event *dto.Event, rErr error) {
				ms.Suite.Require().Nil(event)
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
				ms.Suite.Require().Len(logErr.Code(), model.ErrEventDateBusyCode)
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
			tc.checkResponseBody(event, err)
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
		Date:        dateInFuture.String(),
		Duration:    "45m",
		Description: &desc,
		NotifyTerm:  &nfTermOk,
	})
	ms.Suite.Require().NotNil(newEvent)
	ms.eventIds = append(ms.eventIds, newEvent.ID)
	ms.Suite.Require().NoError(err)

	testCases := []struct {
		name              string
		ID                string
		checkResponseBody func(event *dto.Event, err error)
	}{
		{
			name: "exists event",
			ID:   newEvent.ID,
			checkResponseBody: func(event *dto.Event, rErr error) {
				ms.Suite.Require().NotNil(event)
				ms.Suite.Require().NoError(rErr)
				ms.Suite.Require().Equal(newEvent.ID, event.ID)
			},
		}, {
			name: "wrong event ID",
			ID:   "ascascas",
			checkResponseBody: func(event *dto.Event, rErr error) {
				ms.Suite.Require().NotNil(event)
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
			},
		}, {
			name: "not exists event",
			ID:   uuid.New().String(),
			checkResponseBody: func(event *dto.Event, rErr error) {
				ms.Suite.Require().NotNil(event)
				invErr := errx.NotFound{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
			},
		},
	}

	for _, tc := range testCases {
		ms.Suite.Run(tc.name, func() {
			tc.checkResponseBody(ms.client.GetByID(ctx, tc.ID))
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
			Date:     dateInFuture.String(),
			Duration: "60m",
		}, {
			Title:    "Test event 2",
			Date:     dateInFuture.Add(time.Hour * 2).String(),
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
	dateOk := dateInFuture.Add(time.Hour * 4).String()
	durOk := "90m"
	desc := "new description"
	nfTermOk := "3h"

	dateDup := dateInFuture.Add(time.Minute * 30 * -1).String()
	durDup := "60m"

	titleErr := ""
	dateErr := "-34t---"
	durErr := "45o"
	nfTermErr := "6o"

	testCases := []struct {
		name              string
		inputUpdate       dto.EventUpdate
		ID                string
		checkResponseBody func(error)
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
			checkResponseBody: func(rErr error) {
				ms.Suite.Require().NoError(rErr)
				event, err := ms.client.GetByID(ctx, events[1].ID)
				ms.Suite.Require().NotNil(event)
				ms.Suite.Require().NoError(err)

				// убедиться что событие обновилось
				ms.Suite.Require().Equal(events[1].ID, event.ID)
				ms.Suite.Require().Equal(titleOk, event.Title)
				ms.Suite.Require().Equal(dateOk, event.Date.String())
				ms.Suite.Require().Equal(durOk, event.Duration)
				ms.Suite.Require().Equal(desc, event.Description)
				ms.Suite.Require().Equal(nfTermOk, event.NotifyTerm)
			},
		}, {
			name: "event move to occupied date",
			inputUpdate: dto.EventUpdate{
				Date:     &dateDup,
				Duration: &durDup,
			},
			ID: events[1].ID,
			checkResponseBody: func(rErr error) {
				logErr := errx.Logic{}
				ms.Suite.Require().ErrorAs(rErr, &logErr)
				ms.Suite.Require().Len(logErr.Code(), model.ErrEventDateBusyCode)
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
			checkResponseBody: func(rErr error) {
				invErr := errx.Invalid{}
				ms.Suite.Require().ErrorAs(rErr, &invErr)
				ms.Suite.Require().Len(invErr.Errors, 3)
			},
		},
	}

	for _, tc := range testCases {
		ms.Suite.Run(tc.name, func() {
			tc.checkResponseBody(ms.client.Update(ctx, tc.ID, tc.inputUpdate))
		})
	}
}

func (ms *MainSuiteTest) TestDelete() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dateInFuture := time.Now().Add(time.Hour * 24 * 30)
	newEvent, err := ms.client.Create(ctx, dto.EventCreate{
		Title:    "Test event",
		Date:     dateInFuture.String(),
		Duration: "45m",
	})
	ms.Suite.Require().NotNil(newEvent)
	ms.Suite.Require().NoError(err)

	ms.Suite.Run("removing exists event", func() {
		var event *dto.Event

		err = ms.client.Delete(ctx, newEvent.ID)
		ms.Suite.Require().NoError(err)

		event, err = ms.client.GetByID(ctx, newEvent.ID)
		ms.Suite.Require().Nil(event)
		rErr := errx.NotFound{}
		ms.Suite.Require().ErrorAs(err, &rErr)
	})

	ms.Suite.Run("removing not exists event", func() {
		err = ms.client.Delete(ctx, newEvent.ID)
		ms.Suite.Require().NoError(err)
		rErr := errx.NotFound{}
		ms.Suite.Require().ErrorAs(err, &rErr)
	})
}

/*
func (ms *MainSuiteTest) TestList() {
	events := addEvents(es, [][]byte{
		[]byte(`{
				"title": "Test event 1",
				"date": "2023-02-19T20:00:00.417Z",
				"duration": "45m"
			}`),
		[]byte(`{
				"title": "Test event 2",
				"date": "2023-02-19T09:00:00.417Z",
				"duration": "90m"
			}`),
		[]byte(`{
				"title": "Test event 3",
				"date": "2023-02-13T15:00:00.417Z",
				"duration": "60m"
			}`),
		[]byte(`{
				"title": "Test event",
				"date": "2023-03-10T10:00:00.417Z",
				"duration": "45m"
			}`),
		[]byte(`{
				"title": "Test event",
				"date": "2023-03-22T20:00:00.417Z",
				"duration": "45m",
				"description": "№348239",
				"notifyTerm": "3h"
			}`),
	})

	testCases := []struct {
		name         string
		rangeType    string
		date         string
		expectedCode int
		expectedIDs  []string
	}{
		{
			name:         "day (whole 2023-02-19) events",
			rangeType:    "day",
			date:         "2023-02-19T22:00:00.417Z",
			expectedCode: http.StatusOK,
			expectedIDs: []string{
				events[0].ID,
				events[1].ID,
			},
		}, {
			name:         "month (whole march) events",
			rangeType:    "month",
			date:         "2023-03-19T22:00:00.417Z",
			expectedCode: http.StatusOK,
			expectedIDs: []string{
				events[3].ID,
				events[4].ID,
			},
		}, {
			name:         "week feb 13-19",
			rangeType:    "week",
			date:         "2023-02-15T22:00:00.417Z",
			expectedCode: http.StatusOK,
			expectedIDs: []string{
				events[0].ID,
				events[1].ID,
				events[2].ID,
			},
		}, {
			name:         "empty list",
			rangeType:    "week",
			date:         "2023-02-25T22:00:00.417Z",
			expectedCode: http.StatusOK,
			expectedIDs:  []string{},
		}, {
			name:         "wrong range type",
			rangeType:    "hour",
			date:         "2023-02-25T22:00:00.417Z",
			expectedCode: http.StatusBadRequest,
			expectedIDs:  []string{},
		}, {
			name:         "wrong date",
			rangeType:    "day",
			date:         "-------2023-02-25T22:00:00.417Z-----",
			expectedCode: http.StatusBadRequest,
			expectedIDs:  []string{},
		},
	}
	for _, tc := range testCases {
		es.Suite.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			requestURL := fmt.Sprintf(
				"%s/events/list/%s?date=%s",
				es.testServer.URL,
				tc.rangeType,
				url.QueryEscape(tc.date),
			)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
			es.Suite.Require().NoError(err)
			req.Header.Set("Authorization", ValidUserEmail)

			res, err := http.DefaultClient.Do(req)
			es.Suite.Require().NoError(err)
			defer func() {
				_ = res.Body.Close()
			}()
			es.Suite.Require().Equal(tc.expectedCode, res.StatusCode)
			if tc.expectedCode == http.StatusOK {
				var resp []dto.Event
				actualIDs := make([]string, 0)
				err = json.NewDecoder(res.Body).Decode(&resp)
				es.Suite.Require().NoError(err)
				for _, event := range resp {
					actualIDs = append(actualIDs, event.ID)
				}
				es.Suite.Require().Equal(tc.expectedIDs, actualIDs)
			}
		})
	}
}
*/

func TestCalendar(t *testing.T) {
	suite.Run(t, new(MainSuiteTest))
}
