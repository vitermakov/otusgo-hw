package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/leporo/sqlf"
	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/calendar"
	deps "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps/calendar"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/closer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Calendar struct {
	config   config.Config
	logger   logger.Logger
	deps     *deps.Deps
	services *deps.Services
	closer   *closer.Closer
}

func NewCalendar(config config.Config) App {
	return &Calendar{config: config, closer: closer.NewCloser()}
}

func (ca *Calendar) Initialize(ctx context.Context) error {
	logLevel, err := logger.ParseLevel(ca.config.Logger.Level)
	if err != nil {
		return fmt.Errorf("'%s': %w", ca.config.Logger.Level, err)
	}
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
		ca.closer.Register("DB", closeFn)

		if pool == nil {
			return fmt.Errorf("unable start logger: %w", err)
		}
		dbPool = pool
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
	ca.deps = &deps.Deps{
		Repos:  repos,
		Logger: ca.logger,
		Clock:  clock.New(),
	}

	ca.services = deps.NewServices(ca.deps)

	return nil
}

func (ca *Calendar) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	restServer, closeFn := http.NewHandledServer(ca.config.Servers.HTTP, ca.services, ca.deps)
	ca.closer.Register("REST Server", closeFn)
	grpcServer, closeFn := grpc.NewHandledServer(ca.config.Servers.GRPC, ca.services, ca.deps)
	ca.closer.Register("GRPC Server", closeFn)

	go func() {
		ca.logger.Info("HTTP server starting")
		if err := restServer.Start(); err != nil {
			ca.logger.Error("failed to start HTTP server: %w", err)
			cancel()
		}
	}()
	go func() {
		ca.logger.Info("GRPC server starting")
		if err := grpcServer.Start(); err != nil {
			ca.logger.Error("failed to start GRPC server: %w", err)
			cancel()
		}
	}()

	ca.logger.Info("calendar is running...")
	<-ctx.Done()

	return nil
}

func (ca *Calendar) Close() {
	// 10 секунд на завершение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ca.closer.Close(ctx, ca.logger)
	ca.logger.Info("calendar stopped")
}
