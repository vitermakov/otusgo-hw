package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest"
	"net/http"
	"net/http/httptest"
	"testing"
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
	es.Require().NoError(err)

	// слой репозиториев мы не будем мокать, а будем использовать реализацию memory.
	repos, err := deps.NewRepos(config.Storage{Type: "memory"}, &deps.Resources{})
	es.Require().NoError(err)

	dependencies := &deps.Deps{Repos: repos, Logger: logs}
	services := deps.NewServices(dependencies)
	handlers := NewHandlers(services, logs)

	restServer := rest.NewServer(rest.Config{}, services.Auth, logs)
	handlers.InitRoutes(restServer)

	es.testServer = httptest.NewServer(restServer)

	// для того, чтобы пользователь авторизовался
	_, err = services.User.Add(context.Background(), model.UserCreate{
		Name:  ValidUserEmail,
		Email: ValidUserEmail,
	})
	es.Require().NoError(err)
}

func (es *EventsSuiteTest) TeardownTest() {
	es.testServer.Close()
}

func (es *EventsSuiteTest) TestCreate() {
	testCases := []struct {
		name                  string
		jsonBody              []byte
		expectedCode          int
		expectedRespStatus    string
		expectedRespLogicCode int
	}{
		{
			name: "event add ok",
			jsonBody: []byte(`{
			  "title": "Встреча в Zoom №348239",
			  "date": "2012-02-19T20:00:00.417Z",
			  "duration": 45
			}`),
			expectedCode:          http.StatusOK,
			expectedRespStatus:    "success",
			expectedRespLogicCode: http.StatusOK,
		},
	}

	requestUrl := fmt.Sprintf("%s%s", es.testServer.URL, "/events")
	for _, tc := range testCases {
		es.Run(tc.name, func() {
			var resp ErrorResponseDTO
			req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBuffer(tc.jsonBody))
			es.Require().NoError(err)
			req.Header.Set("Authorization", ValidUserEmail)

			res, err := http.DefaultClient.Do(req)
			es.Require().NoError(err)
			defer func() {
				_ = res.Body.Close()
			}()
			err = json.NewDecoder(res.Body).Decode(&resp)
			es.Require().NoError(err)

			es.Require().Equal(tc.expectedCode, res.StatusCode)
			es.Require().Equal(tc.expectedRespStatus, resp.Status)
			es.Require().Equal(tc.expectedRespLogicCode, resp.Code)
		})
	}
}

func TestEventsApi(t *testing.T) {
	suite.Run(t, new(EventsSuiteTest))
}
