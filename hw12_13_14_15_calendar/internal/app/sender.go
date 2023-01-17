package app

import (
	"context"
	"fmt"
	"time"

	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/sender"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/sender"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/queue"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/closer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/mailer/stdout"
)

type Sender struct {
	config config.Config
	logger logger.Logger
	deps   *deps.Deps
	closer closer.Closer
}

func NewSender(config config.Config) App {
	return &Sender{config: config}
}

func (sa *Sender) Initialize(ctx context.Context) error {
	logLevel, err := logger.ParseLevel(sa.config.Logger.Level)
	if err != nil {
		return fmt.Errorf("'%s': %w", sa.config.Logger.Level, err)
	}
	sa.logger, err = logger.NewLogrus(logger.Config{
		Level:    logLevel,
		FileName: sa.config.Logger.FileName,
	})
	if err != nil {
		return fmt.Errorf("unable start logger: %w", err)
	}

	supportAPI, authFn, err := grpc.NewSupportClient(
		sa.config.API.Calendar.Address,
		sa.config.APILogin,
	)
	if err != nil {
		return fmt.Errorf("error initialize SupportAPI: %w", err)
	}

	listener, closerFn, err := queue.NewConsumer(sa.config.MPQ, sa.logger, sa.config.Notify.QueueListen)
	if err != nil {
		return fmt.Errorf("error start queue listener: %w", err)
	}
	sa.closer.Register("Queue listener", closerFn)

	sa.deps = &deps.Deps{
		Logger:   sa.logger,
		API:      &deps.API{Support: supportAPI},
		APIAuth:  authFn,
		Listener: listener,
		Mailer: stdout.NewMailer(&stdout.Config{
			TmplPath:    sa.config.Mailer.TemplatePath,
			DefaultFrom: sa.config.Mailer.DefaultFrom,
		}),
	}

	return nil
}

func (sa *Sender) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	service := deps.NewSender(
		sa.deps.API.Support, sa.deps.APIAuth, sa.deps.Listener, sa.logger, sa.deps.Mailer,
		sa.config.Notify.QueueListen, sa.config.Mailer.DefaultFrom,
	)

	return service.Run(ctx)
}

func (sa *Sender) Close() {
	// 10 секунд на завершение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	sa.closer.Close(ctx, sa.logger)
	sa.logger.Info("sender stopped")
}
