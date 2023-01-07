package app

import (
	"context"
	"fmt"
	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/scheduler"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/scheduler"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/closer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"time"
)

type Scheduler struct {
	config   config.Config
	logger   logger.Logger
	deps     *deps.Deps
	services *deps.Services
	closer   closer.Closer
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

	//ca.deps = &deps.Deps{Repos: repos, Logger: ca.logger}

	// ca.services = deps.NewServices(ca.deps)

	return nil
}

func (sa *Scheduler) Close() {
	// 10 секунд на завершение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := sa.closer.Close(ctx)
	if err != nil {
		sa.logger.Info("calendar stopped: %s", err.Error())
	} else {
		sa.logger.Info("calendar stopped successfully")
	}
}

func (sa *Scheduler) Run(ctx context.Context) error { //nolint:unparam // will be used
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return nil
}
