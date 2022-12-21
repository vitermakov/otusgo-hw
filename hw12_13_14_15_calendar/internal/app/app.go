package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/leporo/sqlf"
	"github.com/pressly/goose/v3"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest"
	stdlog "log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Application struct {
	config    config.Config
	logger    logger.Logger
	resources *Resources
	deps      *Deps
	services  *Services
}

func (app Application) initialize(ctx context.Context) error {
	var err error

	logLevel, _ := logger.ParseLevel(app.config.Logger.Level)
	app.logger, err = logger.NewLogrus(logger.Config{
		Level:    logLevel,
		FileName: app.config.Logger.FileName,
	})
	if err != nil {
		return fmt.Errorf("unable start logger: %v", err)
	}

	resources := &Resources{}
	if app.config.Storage.Type == "pgsql" {
		pgCfg := app.config.Storage.PgConn
		dsnURL := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(pgCfg.User, pgCfg.Password),
			Host:     pgCfg.Host,
			Path:     "/" + pgCfg.DbName,
			RawQuery: "application_name=" + app.config.ServiceId,
		}
		resources.DbPool, err = sql.Open("pgx", dsnURL.String())
		if err != nil {
			return fmt.Errorf("unable to connect to database: %v", err)
		}
		resources.DbPool.SetConnMaxLifetime(20 * time.Second)

		app.logger.Info(nil, "database connected...")

		// запускаем миграции
		if err = goose.SetDialect("postgres"); err != nil {
			return fmt.Errorf("error init migrations %s", err.Error())
		}
		if err = goose.Up(resources.DbPool, "migrations"); err != nil {
			return fmt.Errorf("error make migrations %s", err.Error())
		}
		app.logger.Info(nil, "migrations OK...")

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
	return nil
}

func (app Application) close() {
	stdlog.Println("closing resources")
	if app.resources.DbPool != nil {
		_ = app.resources.DbPool.Close()
	}
}

func (app Application) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	repos, err := NewRepos(app.config.Storage, app.resources)
	if err != nil {
		app.logger.Fatal(nil, "error init data layer %s", err.Error())
	}
	deps := Deps{Repos: repos, logger: app.logger}

	// TODO: все остальное делаем в ДЗ №13
	services := NewServices(deps)

	restServer := rest.NewServer(rest.Config{}, services.Auth, app.logger)

	// TODO: убрать отсюда
	restServer.GET("/", func(r *http.Request) rest.Response {
		return rest.OK("zer-gud", r.Context().Value("user"))
	})
	restServer.GET("/hello", func(r *http.Request) rest.Response {
		return rest.OK("zer-gud", r.Context().Value("user"))
	})
	restServer.GET("/panic", func(r *http.Request) rest.Response {
		n := 1 / 0
		return rest.OK("unreachable", n)
	})

	var wg sync.WaitGroup

	go func() {
		defer func() {
			wg.Done()
		}()

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err = restServer.Stop(ctx); err != nil {
			app.logger.Error(nil, "failed to stop http server: "+err.Error())
		}
	}()
	go func() {
		wg.Add(1)
		if err = restServer.Start(); err != nil {
			app.logger.Error(nil, "failed to start http server: "+err.Error())
			cancel()
		}
	}()

	app.logger.Info(nil, "calendar is running...")

	wg.Wait()

	return nil
}

func New(config config.Config) *Application {
	return &Application{
		config: config,
	}
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
