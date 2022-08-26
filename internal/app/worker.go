package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hevela/statements/internal/usecase"
)

var errStartTimeEmpty = errors.New("startTime is empty")

type worker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	startAt    string
	interval   time.Duration
	calculator usecase.Calculator
}

type Option func(worker *worker)

func withInterval(interval string) Option {
	return func(worker *worker) {
		inter, err := time.ParseDuration(interval)
		if err != nil {
			log.Fatal(err)
		}
		worker.interval = inter
	}
}

func withStartAt(startAt string) Option {
	return func(worker *worker) {
		if startAt == "" {
			log.Fatal(errStartTimeEmpty)
		}
		worker.startAt = startAt
	}
}

func withContext(ctx context.Context) Option {
	return func(worker *worker) {
		worker.ctx, worker.cancel = context.WithCancel(ctx)
	}
}

func withCalculator(clc usecase.Calculator) Option {
	return func(worker *worker) {
		worker.calculator = clc
	}
}

func newWorker(opts ...Option) *worker {
	w := &worker{
		ctx:        context.Background(),
		cancel:     nil,
		startAt:    "now",
		interval:   0,
		calculator: nil,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *worker) Start() error {
	startingIn, err := w.setup()
	if err != nil {
		return err
	}
	go time.AfterFunc(startingIn, w.start)
	return nil
}

func (w *worker) setup() (startingIn time.Duration, err error) {
	return durationToNextTime(w.startAt)
}

func (w *worker) start() {
	fmt.Println("start calculating")
	go w.calculator.Run()
	ticker := time.NewTicker(w.interval)
	for {
		select {
		case <-ticker.C:
			fmt.Println("start calculating")
			w.calculator.Run()
		case <-w.ctx.Done():
			ticker.Stop()
			return
		}
	}
}
func (w *worker) Stop() {
	w.cancel()
}

func durationToNextTime(clock string) (time.Duration, error) {
	if clock == "" {
		return time.Duration(0), errStartTimeEmpty
	}
	if clock == "now" {
		return time.Duration(0), nil
	}
	now := time.Now()
	startingTime, err := time.ParseInLocation(time.Kitchen, clock, now.Location())
	if err != nil {
		return 0, err
	}
	startingAt := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		startingTime.Hour(),
		startingTime.Minute(),
		startingTime.Second(),
		0,
		now.Location(),
	)
	if now.After(startingTime) {
		startingTime.Add(24 * time.Hour * time.Duration(1))
	}
	timeToStart := startingAt.Sub(now)
	return timeToStart, nil
}
