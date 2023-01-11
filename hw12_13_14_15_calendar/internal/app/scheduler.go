package app

import (
	"context"
	"fmt"
	"time"

	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/scheduler"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/scheduler"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/queue"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

const (
	defNotifyCheckingTime  = time.Minute
	defCleanerCheckingTime = time.Hour * 24
)

type Scheduler struct {
	config config.Config
	logger logger.Logger
	deps   *deps.Deps
}

func NewScheduler(config config.Config) App {
	return &Scheduler{config: config}
}

func (sa *Scheduler) Initialize(ctx context.Context) error {
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

	publisher, err := queue.NewProducer(sa.config.MPQ, sa.logger, sa.config.Notify.QueuePublish)
	if err != nil {
		return fmt.Errorf("error queue publisher: %w", err)
	}

	sa.deps = &deps.Deps{
		API:       &deps.API{Support: supportAPI},
		Logger:    sa.logger,
		Publisher: publisher,
	}
	return nil
}

func (sa *Scheduler) Close() {
	sa.logger.Info("scheduler stopped successfully")
}

func (sa *Scheduler) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	supAPI := sa.deps.API.Support

	ct, err := time.ParseDuration(sa.config.Notify.CheckingTime)
	if err != nil {
		ct = defNotifyCheckingTime
		sa.logger.Warn(
			"wrong notifier checkingTime config value '%s', set default '%s'",
			sa.config.Notify.CheckingTime, ct.String(),
		)
	}
	notifier := deps.NewNotifier(supAPI, sa.deps.Publisher, sa.logger, sa.config.Notify.QueuePublish)
	notifierRun := deps.NewRepeated(notifier, ct, sa.logger)
	notifierRun.Repeat(ctx)

	ct, err = time.ParseDuration(sa.config.Cleanup.CheckingTime)
	if err != nil {
		ct = defCleanerCheckingTime
		sa.logger.Warn(
			"wrong cleaner checkingTime config value '%s', set default '%s'",
			sa.config.Cleanup.CheckingTime, ct.String(),
		)
	}
	cleaner := deps.NewCleaner(supAPI, sa.logger, sa.config.Cleanup.StoreTime)
	cleanerRun := deps.NewRepeated(cleaner, ct, sa.logger)
	cleanerRun.Repeat(ctx)

	<-ctx.Done()
	return nil
}
