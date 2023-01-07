package app

import (
	"context"
	_ "github.com/jackc/pgx/v4/stdlib" // pgx driver for database/sql
	stdlog "log"
)

type App interface {
	Initialize(ctx context.Context) error
	Run(ctx context.Context) error
	Close()
}

// Execute шаблонная функция выполнения приложения
func Execute(ctx context.Context, app App) {
	var err error

	// пропишем defer на закрытие приложения до инициализации
	defer app.Close()

	err = app.Initialize(ctx)
	if err != nil {
		stdlog.Fatalf("не удалось инициализировать приложение: %s", err.Error())
	}
	err = app.Run(ctx)
	if err != nil {
		stdlog.Fatalf("не удалось запустить приложение: %s", err.Error())
	}
}
