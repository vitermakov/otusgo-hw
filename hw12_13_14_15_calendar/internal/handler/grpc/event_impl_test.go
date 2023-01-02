package grpc

import (
	"context"
	"fmt"
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

func (es *EventsSuiteTest) TeardownTest() {
	fmt.Println(3333)
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

			meta := metadata.New(nil)
			meta.Append("authorization", ValidUserEmail)

			ctx = metadata.NewOutgoingContext(ctx, meta)
			event, err := es.evClient.Create(ctx, tc.inputCreate)

			e, ok := status.FromError(err)
			es.True(ok, "error is not status")
			es.Equal(tc.expectedCode, e.Code())
			if e.Code() == codes.OK {
				es.NotEmpty(event)
			}
		})
	}
}

func TestEventsApi(t *testing.T) {
	suite.Run(t, new(EventsSuiteTest))
}
