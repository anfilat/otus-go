package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/server/httpserver"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage/initstorage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()

	if isVersionCommand() {
		printVersion()
		os.Exit(0)
	}

	mainCtx, cancel := context.WithCancel(context.Background())

	go watchSignals(cancel)

	config, err := newConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level, nil, config.Logger.File)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("starting calendar")

	db, err := initstorage.New(mainCtx, config.Database.Inmem, config.Database.Connect)
	if err != nil {
		logg.Fatal(err)
	}

	calendar := app.New(logg, db)

	server := httpserver.NewServer(calendar, logg)
	go func() {
		err := server.Start(config.HTTP.Host + ":" + config.HTTP.Port)
		if err != nil {
			logg.Error(err)
			cancel()
		}
	}()

	logg.Info("calendar is running...")

	<-mainCtx.Done()

	logg.Info("stopping calendar")
	cancel()
	shutDown(logg, server, db)
	logg.Info("calendar is stopped")
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}

func shutDown(logg logger.Logger, server httpserver.Server, db storage.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logg.Error(err)
	}
	if err := db.Close(ctx); err != nil {
		logg.Error(err)
	}
}
