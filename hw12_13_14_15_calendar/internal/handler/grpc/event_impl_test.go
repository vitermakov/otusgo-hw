package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers"
	grpcServ "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"strconv"
	"testing"
	"time"
)

const (
	ValidUserEmail = "auth@otus.ru"
)

type EventsSuiteTest struct {
	suite.Suite
	grpcServer *grpcServ.Server
	conn       *grpc.ClientConn
	evClient   events.EventsClient
}

func (es *EventsSuiteTest) SetupTest() {
	cfg := servers.Config{
		Host: "127.0.0.1",
		Port: 50051,
	}
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
	es.grpcServer = grpcServ.NewServer(cfg, services.Auth, logs)
	InitHandlers(es.grpcServer, services, dependencies)

	go func() {
		err := es.grpcServer.Start()
		es.Require().NoError(err)
	}()
	// сервер запускается не сразу
	for i := 0; i < 10; i++ {
		<-time.After(time.Millisecond * 10)
		es.conn, err = grpc.Dial(
			net.JoinHostPort(cfg.GetHost(), strconv.Itoa(cfg.GetPort())),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			break
		}
	}
	es.Require().NoError(err)
	es.evClient = events.NewEventsClient(es.conn)
	// для того, чтобы пользователь авторизовался
	_, err = services.User.Add(context.Background(), model.UserCreate{
		Name:  ValidUserEmail,
		Email: ValidUserEmail,
	})
	es.Require().NoError(err)
}

func (es *EventsSuiteTest) TearDownTest() {
	_ = es.conn.Close()
	es.grpcServer.Stop()
}

func (es *EventsSuiteTest) TestCreate() {
	dateOk, _ := time.Parse(time.RFC3339, "2023-02-19T20:00:00.417Z")
	descOk := "№348239"
	termOk := time.Hour * 24 * 180
	termWrong := time.Hour * 24 * -1
	dateDup, _ := time.Parse(time.RFC3339, "2023-02-19T19:30:00.417Z")

	testCases := []struct {
		name         string
		inputCreate  *events.CreateEvent
		expectedCode codes.Code
	}{
		{
			name: "event add ok",
			inputCreate: &events.CreateEvent{
				Title:       "Встреча в Zoom",
				Date:        timestamppb.New(dateOk),
				Duration:    durationpb.New(time.Minute * 120),
				Description: &descOk,
				NotifyTerm:  durationpb.New(termOk),
			},
			expectedCode: codes.OK,
		}, {
			name:         "event empty",
			inputCreate:  nil,
			expectedCode: codes.Internal,
		}, {
			name: "event wrong data",
			inputCreate: &events.CreateEvent{
				Title:      "Встреча в Zoom",
				Date:       timestamppb.New(time.Time{}),
				Duration:   durationpb.New(time.Minute * 120 * -1),
				NotifyTerm: durationpb.New(termWrong),
			},
			expectedCode: codes.InvalidArgument,
		}, {
			name: "duplicate entry",
			inputCreate: &events.CreateEvent{
				Title:    "Встреча в Zoom №348239",
				Date:     timestamppb.New(dateDup),
				Duration: durationpb.New(time.Minute * 50),
			},
			expectedCode: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		es.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			event, err := es.evClient.Create(auth(es, ctx), tc.inputCreate)
			e, ok := status.FromError(err)
			es.True(ok, "error is not status")
			es.Equal(tc.expectedCode, e.Code())
			if e.Code() == codes.OK {
				es.NotEmpty(event)
			}
		})
	}
}

func (es *EventsSuiteTest) TestGetByID() {
	dateOk, _ := time.Parse(time.RFC3339, "2023-02-19T20:00:00.417Z")
	descOk := "№348239"
	termOk := time.Hour * 24 * 180
	items := addEvents(es, []events.CreateEvent{
		{
			Title:       "Test event",
			Date:        timestamppb.New(dateOk),
			Duration:    durationpb.New(time.Minute * 45),
			Description: &descOk,
			NotifyTerm:  durationpb.New(termOk),
		},
	})
	testCases := []struct {
		name         string
		ID           string
		expectedCode codes.Code
	}{
		{
			name:         "exists event",
			ID:           items[0].ID,
			expectedCode: codes.OK,
		}, {
			name:         "wrong event ID",
			ID:           "ascascas",
			expectedCode: codes.NotFound,
		}, {
			name:         "not exists event",
			ID:           uuid.New().String(),
			expectedCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		es.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			event, err := es.evClient.GetByID(auth(es, ctx), &events.EventIDReq{ID: tc.ID})
			e, ok := status.FromError(err)
			es.True(ok, "error is not status")
			es.Equal(tc.expectedCode, e.Code())
			if e.Code() == codes.OK {
				es.Equal(tc.ID, event.ID)
			}
		})
	}
}

func (es *EventsSuiteTest) TestUpdate() {
	date1, _ := time.Parse(time.RFC3339, "2023-02-19T20:00:00.417Z")
	desc1 := "№665"
	date2, _ := time.Parse(time.RFC3339, "2023-02-19T10:00:00.417Z")
	desc2 := "№666"
	items := addEvents(es, []events.CreateEvent{
		{
			Title:       "Test event",
			Date:        timestamppb.New(date1),
			Duration:    durationpb.New(time.Minute * 45),
			Description: &desc1,
		}, {
			Title:       "Test event",
			Date:        timestamppb.New(date2),
			Duration:    durationpb.New(time.Minute * 60),
			Description: &desc2,
		},
	})
	titleOk := "Test event 1 updated"
	dateOK, _ := time.Parse(time.RFC3339, "2023-02-19T19:00:00.417Z")
	dateDup, _ := time.Parse(time.RFC3339, "2023-02-19T10:30:00.417Z")
	titleWrong := ""
	testCases := []struct {
		name         string
		inputUpdate  *events.UpdateEvent
		expectedCode codes.Code
	}{
		{
			name: "event update ok",
			inputUpdate: &events.UpdateEvent{
				ID:       items[0].ID,
				Title:    &titleOk,
				Date:     timestamppb.New(dateOK),
				Duration: durationpb.New(time.Minute * 60),
			},
			expectedCode: codes.OK,
		}, {
			name: "event move to occupied date",
			inputUpdate: &events.UpdateEvent{
				ID:       items[0].ID,
				Date:     timestamppb.New(dateDup),
				Duration: durationpb.New(time.Minute * 60),
			},
			expectedCode: codes.InvalidArgument,
		}, {
			name: "event wrong data",
			inputUpdate: &events.UpdateEvent{
				ID:         items[1].ID,
				Title:      &titleWrong,
				Date:       timestamppb.New(time.Time{}),
				Duration:   durationpb.New(time.Minute * 60 * -1),
				NotifyTerm: durationpb.New(time.Hour * 24 * 160 * -1),
			},
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		es.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			_, err := es.evClient.Update(auth(es, ctx), tc.inputUpdate)
			e, ok := status.FromError(err)
			es.True(ok, "error is not status")
			es.Equal(tc.expectedCode, e.Code())
		})
	}
}

func (es *EventsSuiteTest) TestDelete() {
	dateOk, _ := time.Parse(time.RFC3339, "2023-02-19T20:00:00.417Z")
	descOk := "№348239"
	termOk := time.Hour * 24 * 180
	items := addEvents(es, []events.CreateEvent{
		{
			Title:       "Test event",
			Date:        timestamppb.New(dateOk),
			Duration:    durationpb.New(time.Minute * 45),
			Description: &descOk,
			NotifyTerm:  durationpb.New(termOk),
		},
	})
	es.Run("removing exists event", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err := es.evClient.Delete(auth(es, ctx), &events.EventIDReq{ID: items[0].ID})
		e, ok := status.FromError(err)
		es.True(ok, "error is not status")
		es.Equal(codes.OK, e.Code())

		// убедиться что события нет
		_, err = es.evClient.GetByID(auth(es, ctx), &events.EventIDReq{ID: items[0].ID})
		e, ok = status.FromError(err)
		es.True(ok, "error is not status")
		es.Equal(codes.NotFound, e.Code())
	})

	es.Run("removing exists event", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err := es.evClient.Delete(auth(es, ctx), &events.EventIDReq{ID: uuid.New().String()})
		e, ok := status.FromError(err)
		es.True(ok, "error is not status")
		es.Equal(codes.NotFound, e.Code())
	})
}

func (es *EventsSuiteTest) TestList() {
	date1, _ := time.Parse(time.RFC3339, "2023-02-19T20:00:00.417Z")
	date2, _ := time.Parse(time.RFC3339, "2023-02-19T09:00:00.417Z")
	date3, _ := time.Parse(time.RFC3339, "2023-02-13T15:00:00.417Z")
	date4, _ := time.Parse(time.RFC3339, "2023-03-10T10:00:00.417Z")
	date5, _ := time.Parse(time.RFC3339, "2023-03-22T20:00:00.417Z")

	items := addEvents(es, []events.CreateEvent{
		{
			Title:    "Test event 1",
			Date:     timestamppb.New(date1),
			Duration: durationpb.New(time.Minute * 45),
		}, {
			Title:    "Test event 2",
			Date:     timestamppb.New(date2),
			Duration: durationpb.New(time.Minute * 90),
		}, {
			Title:    "Test event 3",
			Date:     timestamppb.New(date3),
			Duration: durationpb.New(time.Minute * 60),
		}, {
			Title:    "Test event 4",
			Date:     timestamppb.New(date4),
			Duration: durationpb.New(time.Minute * 45),
		}, {
			Title:    "Test event 5",
			Date:     timestamppb.New(date5),
			Duration: durationpb.New(time.Minute * 45),
		},
	})

	date1, _ = time.Parse(time.RFC3339, "2023-02-19T22:00:00.417Z")
	date2, _ = time.Parse(time.RFC3339, "2023-03-19T22:00:00.417Z")
	date3, _ = time.Parse(time.RFC3339, "2023-02-15T22:00:00.417Z")
	date4, _ = time.Parse(time.RFC3339, "2023-02-25T22:00:00.417Z")
	date5, _ = time.Parse(time.RFC3339, "2023-02-25T22:00:00.417Z")

	testCases := []struct {
		name         string
		rangeType    events.RangeType
		date         *timestamppb.Timestamp
		expectedCode codes.Code
		expectedIDs  []string
	}{
		{
			name:         "day (whole 2023-02-19) events",
			rangeType:    events.RangeType_RANGE_TYPE_DAY,
			date:         timestamppb.New(date1),
			expectedCode: codes.OK,
			expectedIDs: []string{
				items[0].ID,
				items[1].ID,
			},
		}, {
			name:         "month (whole march) events",
			rangeType:    events.RangeType_RANGE_TYPE_MONTH,
			date:         timestamppb.New(date2),
			expectedCode: codes.OK,
			expectedIDs: []string{
				items[3].ID,
				items[4].ID,
			},
		}, {
			name:         "week feb 13-19",
			rangeType:    events.RangeType_RANGE_TYPE_WEEK,
			date:         timestamppb.New(date3),
			expectedCode: codes.OK,
			expectedIDs: []string{
				items[0].ID,
				items[1].ID,
				items[2].ID,
			},
		}, {
			name:         "empty list",
			rangeType:    events.RangeType_RANGE_TYPE_WEEK,
			date:         timestamppb.New(date4),
			expectedCode: codes.OK,
			expectedIDs:  []string{},
		}, {
			name:         "wrong range type",
			rangeType:    100,
			date:         timestamppb.New(date5),
			expectedCode: codes.InvalidArgument,
			expectedIDs:  []string{},
		}, {
			name:         "wrong date",
			rangeType:    events.RangeType_RANGE_TYPE_DAY,
			date:         timestamppb.New(time.Time{}),
			expectedCode: codes.InvalidArgument,
			expectedIDs:  []string{},
		},
	}
	for _, tc := range testCases {
		es.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			items, err := es.evClient.GetListOnDate(auth(es, ctx), &events.ListOnDateReq{
				Date:      tc.date,
				RangeType: tc.rangeType,
			})
			e, ok := status.FromError(err)
			es.True(ok, "error is not status")
			es.Equal(tc.expectedCode, e.Code())

			if tc.expectedCode == codes.OK {
				actualIDs := make([]string, 0)
				for _, event := range items.List {
					actualIDs = append(actualIDs, event.ID)
				}
				es.Require().Equal(tc.expectedIDs, actualIDs)
			}
		})
	}
}

func TestEventsApi(t *testing.T) {
	suite.Run(t, new(EventsSuiteTest))
}

func addEvents(es *EventsSuiteTest, items []events.CreateEvent) []*events.Event {
	es.T().Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result := make([]*events.Event, len(items))

	meta := metadata.New(nil)
	meta.Append("authorization", ValidUserEmail)

	ctx = metadata.NewOutgoingContext(ctx, meta)
	for i, inputCreate := range items {
		event, err := es.evClient.Create(ctx, &inputCreate)
		e, ok := status.FromError(err)
		es.True(ok, "error is not status")
		es.Equal(codes.OK, e.Code())
		es.NotEmpty(event)

		result[i] = event
	}

	return result
}
func auth(es *EventsSuiteTest, ctx context.Context) context.Context {
	es.T().Helper()
	meta := metadata.New(nil)
	meta.Append("authorization", ValidUserEmail)
	return metadata.NewOutgoingContext(ctx, meta)
}
