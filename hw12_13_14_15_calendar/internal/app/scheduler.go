package app

import (
	"context"
	"fmt"
	"time"

	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/scheduler"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/scheduler"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/closer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Scheduler struct {
	config config.Config
	logger logger.Logger
	deps   *deps.Deps
	closer *closer.Closer
}

func NewScheduler(config config.Config) App {
	return &Scheduler{config: config, closer: closer.NewCloser()}
}

func (sa *Scheduler) Initialize(ctx context.Context) error {
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

	publisher, closerFn, err := queue.NewProducer(sa.config.MPQ, sa.logger, sa.config.Notify.QueuePublish)
	if err != nil {
		return fmt.Errorf("error start queue publisher: %w", err)
	}
	sa.closer.Register("Queue publisher", closerFn)

	sa.deps = &deps.Deps{
		API:       &deps.API{Support: supportAPI},
		APIAuth:   authFn,
		Logger:    sa.logger,
		Publisher: publisher,
	}
	return nil
}

func (sa *Scheduler) Close() {
	// 10 секунд на завершение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	sa.closer.Close(ctx, sa.logger)
	sa.logger.Info("scheduler stopped")
}

func (sa *Scheduler) Run(ctx context.Context) error {
	supAPI := sa.deps.API.Support

	notifier := deps.NewNotifier(supAPI, sa.deps.APIAuth, sa.deps.Publisher, sa.logger, sa.config.Notify.QueuePublish)
	notifierRun := deps.NewRepeated(notifier, sa.config.Notify.CheckingTime, sa.logger)
	notifierRun.Repeat(ctx)

	cleaner := deps.NewCleaner(supAPI, sa.deps.APIAuth, sa.logger, sa.config.Cleanup.StoreTime)
	cleanerRun := deps.NewRepeated(cleaner, sa.config.Cleanup.CheckingTime, sa.logger)
	cleanerRun.Repeat(ctx)

	<-ctx.Done()
	return nil
}
