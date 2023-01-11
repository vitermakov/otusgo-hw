package app

import (
	"context"
	"fmt"

	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/sender"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/sender"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/queue"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/mailer/stdout"
)

type Sender struct {
	config config.Config
	logger logger.Logger
	deps   *deps.Deps
}

func NewSender(config config.Config) App {
	return &Sender{config: config}
}

func (sa *Sender) Initialize(ctx context.Context) error {
	var err error
	logLevel, _ := logger.ParseLevel(sa.config.Logger.Level)
	sa.logger, err = logger.NewLogrus(logger.Config{
		Level:    logLevel,
		FileName: sa.config.Logger.FileName,
	})
	if err != nil {
		return fmt.Errorf("unable start logger: %w", err)
	}

	supportAPI, err := grpc.NewSupportClient(sa.config.API.Calendar.Address)
	if err != nil {
		return fmt.Errorf("error initialize SupportAPI: %w", err)
	}

	listener, err := queue.NewConsumer(sa.config.MPQ, sa.logger, sa.config.Notify.QueueListen)
	if err != nil {
		return fmt.Errorf("error queue listener: %w", err)
	}

	sa.deps = &deps.Deps{
		Logger:   sa.logger,
		API:      &deps.API{Support: supportAPI},
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
		sa.deps.API.Support, sa.deps.Listener, sa.logger, sa.deps.Mailer,
		sa.config.Notify.QueueListen, sa.config.Mailer.DefaultFrom,
	)

	return service.Run(ctx)
}

func (sa *Sender) Close() {
	sa.logger.Info("sender stopped successfully")
}
