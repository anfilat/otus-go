package app

import (
	"context"
	"time"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type app struct {
	logger  logger.Logger
	storage storage.Storage
}

func (a *app) Create(ctx context.Context, userID int, title, desc string, start, stop time.Time, notif *time.Duration) (id int, err error) {
	if userID == 0 {
		err = ErrNoUserID
		return
	}
	if title == "" {
		err = ErrEmptyTitle
		return
	}
	if start.After(stop) {
		start, stop = stop, start
	}
	if time.Now().After(start) {
		err = ErrStartInPast
		return
	}
	isBusy, err := a.storage.IsTimeBusy(ctx, start, stop, 0)
	if err != nil {
		return
	}
	if isBusy {
		err = ErrDateBusy
		return
	}

	return a.storage.Create(ctx, storage.Event{
		Title:        title,
		Start:        start,
		Stop:         stop,
		Description:  desc,
		UserID:       userID,
		Notification: notif,
	})
}

func (a *app) Update(ctx context.Context, id int, change storage.Event) error {
	if change.Title == "" {
		return ErrEmptyTitle
	}
	if change.Start.After(change.Stop) {
		change.Start, change.Stop = change.Stop, change.Start
	}
	if time.Now().After(change.Start) {
		return ErrStartInPast
	}
	isBusy, err := a.storage.IsTimeBusy(ctx, change.Start, change.Stop, id)
	if err != nil {
		return err
	}
	if isBusy {
		return ErrDateBusy
	}

	return a.storage.Update(ctx, id, change)
}

func (a *app) Delete(ctx context.Context, id int) error {
	return a.storage.Delete(ctx, id)
}

func (a *app) ListAll(ctx context.Context) ([]storage.Event, error) {
	return a.storage.ListAll(ctx)
}

func (a *app) ListDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListDay(ctx, date)
}

func (a *app) ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListWeek(ctx, date)
}

func (a *app) ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListMonth(ctx, date)
}
