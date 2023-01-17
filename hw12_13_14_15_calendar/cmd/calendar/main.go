package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app"
	config "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config/calendar"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	app.Execute(ctx, app.NewCalendar(cfg))
}
