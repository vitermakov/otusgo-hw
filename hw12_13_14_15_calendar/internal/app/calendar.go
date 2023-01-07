package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/leporo/sqlf"
	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/calendar"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/closer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"sync"
	"time"
)

type Calendar struct {
	config   config.Config
	logger   logger.Logger
	deps     *deps.Deps
	services *deps.Services
	closer   closer.Closer
}

func NewCalendar(config config.Config) App {
	return &Calendar{config: config}
}

func (ca *Calendar) Initialize(ctx context.Context) error {
	var err error
	logLevel, _ := logger.ParseLevel(ca.config.Logger.Level)
	ca.logger, err = logger.NewLogrus(logger.Config{
		Level:    logLevel,
		FileName: ca.config.Logger.FileName,
	})
	if err != nil {
		return fmt.Errorf("unable start logger: %w", err)
	}

	var dbPool *sql.DB
	if ca.config.Storage.Type == "pgsql" {
		pool, closeFn := deps.NewPgConn(ca.config.ServiceID, ca.config.Storage.PGConn, ca.logger)
		ca.closer.Register(closeFn)

		if pool == nil {
			return fmt.Errorf("unable start logger: %w", err)
		}
		// устанавливаем диалект билдера запросов
		sqlf.SetDialect(sqlf.PostgreSQL)
		// это костыль, так как при большом количестве запросов он подтекает
		go func() {
			for {
				sqlf.PostgreSQL.ClearCache()
				sqlf.NoDialect.ClearCache()
				select {
				case <-ctx.Done():
					return
				case <-time.After(30 * time.Minute):
				}
			}
		}()
	}

	repos, err := deps.NewRepos(ca.config.Storage, dbPool)
	if err != nil {
		return fmt.Errorf("error init data layer %w", err)
	}
	ca.deps = &deps.Deps{Repos: repos, Logger: ca.logger}

	ca.services = deps.NewServices(ca.deps)

	return nil
}

func (ca *Calendar) Close() {
	// 10 секунд на завершение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := ca.closer.Close(ctx)
	if err != nil {
		ca.logger.Info("calendar stopped: %s", err.Error())
	} else {
		ca.logger.Info("calendar stopped successfully")
	}
}

func (ca *Calendar) Run(ctx context.Context) error { //nolint:unparam // will be used
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	restServer := http.NewHandledServer(ca.config.Servers.HTTP, ca.services, ca.deps)
	grpcServer := grpc.NewHandledServer(ca.config.Servers.GRPC, ca.services, ca.deps)

	var wg sync.WaitGroup

	// TODO: выделить отдельный механизм для завершающих процедур, убрать в close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		restServer.Stop(ctx)
	}()

	go func() {
		if err := restServer.Start(); err != nil {
			cancel()
		}
	}()

	// TODO: выделить отдельный механизм для завершающих процедур, убрать в close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		grpcServer.Stop()
	}()

	go func() {
		if err := grpcServer.Start(); err != nil {
			cancel()
		}
	}()

	ca.logger.Info("calendar is running...")

	wg.Wait()

	return nil
}
