package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/deps"
	handler "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http"
	stdlog "log"
	"net/url"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // pgx driver for database/sql
	"github.com/leporo/sqlf"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest"
)

type Application struct {
	config    config.Config
	logger    logger.Logger
	resources *deps.Resources
	deps      *deps.Deps
	services  *deps.Services
}

func (app *Application) initialize(ctx context.Context) error {
	var err error
	logLevel, _ := logger.ParseLevel(app.config.Logger.Level)
	app.logger, err = logger.NewLogrus(logger.Config{
		Level:    logLevel,
		FileName: app.config.Logger.FileName,
	})
	if err != nil {
		return fmt.Errorf("unable start logger: %w", err)
	}

	app.resources = &deps.Resources{}
	if app.config.Storage.Type == "pgsql" {
		pgCfg := app.config.Storage.PGConn
		dsnURL := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(pgCfg.User, pgCfg.Password),
			Host:     pgCfg.Host,
			Path:     "/" + pgCfg.DBName,
			RawQuery: "application_name=" + app.config.ServiceID,
		}
		app.resources.DBPool, err = sql.Open("pgx", dsnURL.String())
		if err != nil {
			return fmt.Errorf("unable to connect to database: %w", err)
		}
		app.resources.DBPool.SetConnMaxLifetime(20 * time.Second)

		app.logger.Info("database connected...")

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

	repos, err := deps.NewRepos(app.config.Storage, app.resources)
	if err != nil {
		return fmt.Errorf("error init data layer %w", err)
	}
	app.deps = &deps.Deps{Repos: repos, Logger: app.logger}

	// TODO: все остальное делаем в ДЗ №13.
	app.services = deps.NewServices(app.deps)

	return nil
}

func (app *Application) close() {
	if app.logger != nil {
		app.logger.Info("closing resources")
	} else {
		stdlog.Println("closing resources")
	}
	if app.resources.DBPool != nil {
		_ = app.resources.DBPool.Close()
	}
}

func (app *Application) run(ctx context.Context) error { //nolint:unparam // will be used
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	restServer := rest.NewServer(rest.Config{}, app.services.Auth, app.logger)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := restServer.Stop(ctx); err != nil {
			app.logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		handlers := handler.NewHandlers(app.services, app.deps.Logger)
		handlers.InitRoutes(restServer)
		if err := restServer.Start(); err != nil {
			app.logger.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	app.logger.Info("calendar is running...")

	wg.Wait()

	return nil
}

func (app *Application) Main(ctx context.Context) {
	var err error

	// пропишем defer на закрытие приложения до инициализации
	defer app.close()

	err = app.initialize(ctx)
	if err != nil {
		stdlog.Fatalf("не удалось инициализировать приложение: %s", err.Error())
	}
	err = app.run(ctx)
	if err != nil {
		stdlog.Fatalf("не удалось запустить приложение: %s", err.Error())
	}
}

func New(config config.Config) *Application {
	return &Application{
		config: config,
	}
}
