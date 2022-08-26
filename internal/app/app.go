package app

import (
	"context"
	"github.com/hevela/statements/config"
	"github.com/hevela/statements/internal/usecase"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"path"
)

func Run(cfg *config.Config) {
	// configure statements Calculator
	p, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	p = path.Join(p, cfg.FilesDir)
	clc := usecase.NewCalculator(
		usecase.WithDirPath(p),
		usecase.WithAPIKey(cfg.SendGridAPIKey),
		usecase.WithTemplateID(cfg.TemplateID),
	)
	// init Worker
	wrkr := newWorker(
		withInterval(cfg.Interval),
		withStartAt(cfg.StartAt),
		withContext(context.Background()),
		withCalculator(clc),
	)
	shutdownChan := setupShutdown(wrkr)
	// Run the worker in a goroutine so that it doesn't block
	go func() {
		startWorker(wrkr)
	}()
	<-shutdownChan
	logrus.Info("shutting down")
}

func startWorker(w *worker) {
	logrus.Info("starting worker")
	if err := w.Start(); err != nil {
		logrus.Fatalf("failed to initialize worker: %v", err)
	}
}

func setupShutdown(w *worker) chan struct{} {
	shutdownComplete := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint // we wait...
		w.Stop()
		close(shutdownComplete)
	}()
	return shutdownComplete
}
