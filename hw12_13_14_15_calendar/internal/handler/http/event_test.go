package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type EventsSuiteTest struct {
	suite.Suite
	testServer *httptest.Server
}

type ErrorResponseDTO struct {
	Status  string            `json:"status"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    json.RawMessage   `json:"data"`
	Errors  map[string]string `json:"errors"`
}

const (
	ValidUserEmail = "auth@otus.ru"
)

func (es *EventsSuiteTest) SetupTest() {
	// логгер будет писать в stdout и stderr, а мы будем это перехватывать.
	logs, err := logger.NewLogrus(logger.Config{
		Level: logger.LevelInfo,
	})
	es.Suite.Require().NoError(err)

	// слой репозиториев мы не будем мокать, а будем использовать реализацию memory.
	repos, err := deps.NewRepos(config.Storage{Type: "memory"}, nil)
	es.Suite.Require().NoError(err)

	dependencies := &deps.Deps{Repos: repos, Logger: logs}
	services := deps.NewServices(dependencies)

	restServer := NewHandledServer(config.Server{}, services, dependencies)

	es.testServer = httptest.NewServer(restServer)

	// для того, чтобы пользователь авторизовался
	_, err = services.User.Add(context.Background(), model.UserCreate{
		Name:  ValidUserEmail,
		Email: ValidUserEmail,
	})
	es.Suite.Require().NoError(err)
}

func (es *EventsSuiteTest) TearDownTest() {
	es.testServer.Close()
}

func (es *EventsSuiteTest) TestCreate() {
	testCases := []struct {
		name                  string
		jsonBody              []byte
		expectedCode          int
		expectedRespStatus    string
		expectedRespLogicCode int
		checkResponseBody     func(resp ErrorResponseDTO)
	}{
		{
			name: "event add ok",
			jsonBody: []byte(`{
			  	"title": "Встреча в Zoom",
			  	"date": "2023-02-19T20:00:00.417Z",
			  	"duration": "45m",
				"description": "№348239",
				"notifyTerm": "3h"
			}`),
			expectedCode:          http.StatusOK,
			expectedRespStatus:    "success",
			expectedRespLogicCode: http.StatusOK,
			checkResponseBody: func(resp ErrorResponseDTO) {
				var eventDto dto.Event
				err := json.Unmarshal(resp.Data, &eventDto)
				es.Suite.Require().NoError(err)

				_, err = uuid.Parse(eventDto.ID)
				es.Suite.Require().NoError(err)
				d, _ := time.ParseDuration("45m")
				es.Suite.Require().Equal(d.String(), eventDto.Duration)
				nt, _ := time.ParseDuration("3h")
				es.Suite.Require().Equal(nt.String(), eventDto.NotifyTerm)
			},
		}, {
			name: "crushed json",
			jsonBody: []byte(`{
			  "title": "Встреча в Zoom №348239",
			  "date": "2023-02-19T20:00:00.417Z",
			  "duration": "45m",
			}`),
			expectedCode:          http.StatusBadRequest,
			expectedRespStatus:    "error",
			expectedRespLogicCode: http.StatusBadRequest,
		}, {
			name:                  "event empty",
			jsonBody:              nil,
			expectedCode:          http.StatusUnprocessableEntity,
			expectedRespStatus:    "error",
			expectedRespLogicCode: http.StatusUnprocessableEntity,
			checkResponseBody: func(resp ErrorResponseDTO) {
				es.Suite.Require().Len(resp.Errors, 2)
			},
		}, {
			name: "event wrong data",
			jsonBody: []byte(`{
			  	"title": "Test event",
			  	"date": "----2023-02-19T20:00:00.417Z-----",
			  	"duration": "45q",
				"notifyTerm": "6q"
			}`),
			expectedCode:          http.StatusUnprocessableEntity,
			expectedRespStatus:    "error",
			expectedRespLogicCode: http.StatusUnprocessableEntity,
			checkResponseBody: func(resp ErrorResponseDTO) {
				es.Suite.Require().Len(resp.Errors, 3)
			},
		}, {
			name: "duplicate entry",
			jsonBody: []byte(`{
			  "title": "Встреча в Zoom №348239",
			  "date": "2023-02-19T19:30:00.417Z",
			  "duration": "50m"
			}`),
			expectedCode:          http.StatusBadRequest,
			expectedRespStatus:    "error",
			expectedRespLogicCode: model.ErrEventDateBusyCode,
		},
	}
	requestURL := fmt.Sprintf("%s/events", es.testServer.URL)
	for _, tc := range testCases {
		&deps.Resources{}
		es.Suite.Run(tc.name, func() {
			var resp ErrorResponseDTO
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(tc.jsonBody))
			es.Suite.Require().NoError(err)
			req.Header.Set("Authorization", ValidUserEmail)

			res, err := http.DefaultClient.Do(req)
			es.Suite.Require().NoError(err)
			defer func() {
				_ = res.Body.Close()
			}()
			err = json.NewDecoder(res.Body).Decode(&resp)
			es.Suite.Require().NoError(err)

			es.Suite.Require().Equal(tc.expectedCode, res.StatusCode)
			es.Suite.Require().Equal(tc.expectedRespStatus, resp.Status)
			es.Suite.Require().Equal(tc.expectedRespLogicCode, resp.Code)

			if tc.checkResponseBody != nil {
				tc.checkResponseBody(resp)
			}
		})
	}
}

func (es *EventsSuiteTest) TestGetByID() {
	events := addEvents(es, [][]byte{
		[]byte(`{
			"title": "Test event",
			"date": "2023-02-19T20:00:00.417Z",
			"duration": "45m",
			"description": "№348239",
			"notifyTerm": "3h"
		}`),
	})
	testCases := []struct {
		name         string
		ID           string
		expectedCode int
	}{
		{
			name:         "exists event",
			ID:           events[0].ID,
			expectedCode: http.StatusOK,
		}, {
			name:         "wrong event ID",
			ID:           "ascascas",
			expectedCode: http.StatusBadRequest,
		}, {
			name:         "not exists event",
			ID:           uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		es.Suite.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			requestURL := fmt.Sprintf("%s/events/%s", es.testServer.URL, tc.ID)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
			es.Suite.Require().NoError(err)
			req.Header.Set("Authorization", ValidUserEmail)

			res, err := http.DefaultClient.Do(req)
			es.Suite.Require().NoError(err)
			defer func() {
				_ = res.Body.Close()
			}()

			es.Suite.Require().Equal(tc.expectedCode, res.StatusCode)
		})
	}
}

func (es *EventsSuiteTest) TestUpdate() {
	events := addEvents(es, [][]byte{
		[]byte(`{
			"title": "Test event 1",
			"date": "2023-02-19T20:00:00.417Z",
			"duration": "45m",
			"description": "№665"
		}`),
		[]byte(`{
			"title": "Test event 2",
			"date": "2023-02-19T10:00:00.417Z",
			"duration": "60m",
			"description": "№666"
		}`),
	})
	testCases := []struct {
		name                  string
		jsonBody              []byte
		ID                    string
		expectedCode          int
		expectedRespStatus    string
		expectedRespLogicCode int
		checkResponseBody     func(resp ErrorResponseDTO)
	}{
		{
			name: "event update ok",
			jsonBody: []byte(`{
			  	"title": "Test event 1 updated",
			  	"date": "2023-02-19T19:00:00.417Z",
			  	"duration": "60m",
				"description": "Test description",
				"notifyTerm": "4h"
			}`),
			ID:                    events[0].ID,
			expectedCode:          http.StatusOK,
			expectedRespStatus:    "success",
			expectedRespLogicCode: http.StatusOK,
			checkResponseBody: func(_ ErrorResponseDTO) {
				var eventDto dto.Event
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				// убедиться что событие обновилось
				requestURL := fmt.Sprintf("%s/events/%s", es.testServer.URL, events[0].ID)
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
				es.Suite.Require().NoError(err)
				req.Header.Set("Authorization", ValidUserEmail)

				res, err := http.DefaultClient.Do(req)
				es.Suite.Require().NoError(err)
				defer func() {
					_ = res.Body.Close()
				}()
				es.Suite.Require().Equal(http.StatusOK, res.StatusCode)
				err = json.NewDecoder(res.Body).Decode(&eventDto)
				es.Suite.Require().NoError(err)

				es.Suite.Require().Equal(events[0].ID, eventDto.ID)
				es.Suite.Require().Equal("Test event 1 updated", eventDto.Title)
				d, _ := time.Parse(time.RFC3339, "2023-02-19T19:00:00.417Z")
				es.Suite.Require().Equal(d, eventDto.Date)
				du, _ := time.ParseDuration("60m")
				es.Suite.Require().Equal(du.String(), eventDto.Duration)
				es.Suite.Require().Equal("Test description", eventDto.Description)
				nt, _ := time.ParseDuration("4h")
				es.Suite.Require().Equal(nt.String(), eventDto.NotifyTerm)
			},
		}, {
			name: "event move to occupied date",
			jsonBody: []byte(`{
			  	"date": "2023-02-19T10:30:00.417Z",
			  	"duration": "60m"
			}`),
			ID:                    events[0].ID,
			expectedCode:          http.StatusBadRequest,
			expectedRespStatus:    "error",
			expectedRespLogicCode: model.ErrEventDateBusyCode,
		}, {
			name:                  "crushed json",
			jsonBody:              []byte(`{...}`),
			ID:                    events[1].ID,
			expectedCode:          http.StatusBadRequest,
			expectedRespStatus:    "error",
			expectedRespLogicCode: http.StatusBadRequest,
		}, {
			name: "event wrong data",
			jsonBody: []byte(`{
			  	"title": "",
			  	"date": "-34t---",
			  	"duration": "45o",
				"notifyTerm": "6o"
			}`),
			ID:                    events[1].ID,
			expectedCode:          http.StatusUnprocessableEntity,
			expectedRespStatus:    "error",
			expectedRespLogicCode: http.StatusUnprocessableEntity,
			checkResponseBody: func(resp ErrorResponseDTO) {
				es.Suite.Require().Len(resp.Errors, 3)
			},
		},
	}

	for _, tc := range testCases {
		es.Suite.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			requestURL := fmt.Sprintf("%s/events/%s", es.testServer.URL, tc.ID)

			var resp ErrorResponseDTO
			req, err := http.NewRequestWithContext(ctx, http.MethodPut, requestURL, bytes.NewBuffer(tc.jsonBody))
			es.Suite.Require().NoError(err)
			req.Header.Set("Authorization", ValidUserEmail)

			res, err := http.DefaultClient.Do(req)
			es.Suite.Require().NoError(err)
			defer func() {
				_ = res.Body.Close()
			}()
			err = json.NewDecoder(res.Body).Decode(&resp)
			es.Suite.Require().NoError(err)

			es.Suite.Require().Equal(tc.expectedCode, res.StatusCode)
			es.Suite.Require().Equal(tc.expectedRespStatus, resp.Status)
			es.Suite.Require().Equal(tc.expectedRespLogicCode, resp.Code)

			if tc.checkResponseBody != nil {
				tc.checkResponseBody(resp)
			}
		})
	}
}

func (es *EventsSuiteTest) TestDelete() {
	events := addEvents(es, [][]byte{
		[]byte(`{
				"title": "Test event",
				"date": "2023-02-19T20:00:00.417Z",
				"duration": "45m",
				"description": "№348239",
				"notifyTerm": "3h"
			}`),
	})
	es.Suite.Run("removing exists event", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		requestURL := fmt.Sprintf("%s/events/%s", es.testServer.URL, events[0].ID)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, nil)
		es.Suite.Require().NoError(err)
		req.Header.Set("Authorization", ValidUserEmail)

		res, err := http.DefaultClient.Do(req)
		es.Suite.Require().NoError(err)
		defer func() {
			_ = res.Body.Close()
		}()

		es.Suite.Require().Equal(http.StatusOK, res.StatusCode)

		// убедиться что события нет
		requestURL = fmt.Sprintf("%s/events/%s", es.testServer.URL, events[0].ID)
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
		es.Suite.Require().NoError(err)
		req.Header.Set("Authorization", ValidUserEmail)

		res, err = http.DefaultClient.Do(req)
		es.Suite.Require().NoError(err)
		defer func() {
			_ = res.Body.Close()
		}()

		es.Suite.Require().Equal(http.StatusNotFound, res.StatusCode)
	})

	es.Suite.Run("removing exists event", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		requestURL := fmt.Sprintf("%s/events/%s", es.testServer.URL, uuid.New())
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, nil)
		es.Suite.Require().NoError(err)
		req.Header.Set("Authorization", ValidUserEmail)

		res, err := http.DefaultClient.Do(req)
		es.Suite.Require().NoError(err)
		defer func() {
			_ = res.Body.Close()
		}()

		es.Suite.Require().Equal(http.StatusNotFound, res.StatusCode)
	})
}

func (es *EventsSuiteTest) TestList() {
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

func TestEventsApi(t *testing.T) {
	suite.Run(t, new(EventsSuiteTest))
}

func addEvents(es *EventsSuiteTest, jsonBodies [][]byte) []dto.Event {
	es.Suite.T().Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var resp ErrorResponseDTO
	result := make([]dto.Event, len(jsonBodies))
	requestURL := fmt.Sprintf("%s%s", es.testServer.URL, "/events")
	for i, jsonBody := range jsonBodies {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(jsonBody))
		es.Suite.Require().NoError(err)
		req.Header.Set("Authorization", ValidUserEmail)

		res, err := http.DefaultClient.Do(req)
		es.Suite.Require().NoError(err)
		err = json.NewDecoder(res.Body).Decode(&resp)
		_ = res.Body.Close()
		es.Suite.Require().NoError(err)
		es.Suite.Require().Equal(http.StatusOK, res.StatusCode)

		var eventDto dto.Event
		err = json.Unmarshal(resp.Data, &eventDto)
		es.Suite.Require().NoError(err)

		result[i] = eventDto
	}

	return result
}
