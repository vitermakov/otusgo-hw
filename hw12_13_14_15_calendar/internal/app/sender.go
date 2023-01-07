package app

import (
	"context"
	"fmt"
	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/sender"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/sender"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/closer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"time"
)

type Sender struct {
	config   config.Config
	logger   logger.Logger
	deps     *deps.Deps
	services *deps.Services
	closer   closer.Closer
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

	//ca.deps = &deps.Deps{Repos: repos, Logger: ca.logger}

	// ca.services = deps.NewServices(ca.deps)

	return nil
}

func (sa *Sender) Close() {
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

func (sa *Sender) Run(ctx context.Context) error { //nolint:unparam // will be used
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return nil
}
