package app

import (
	"context"
	"errors"
	"time"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type ListEvents func(ctx context.Context, date time.Time) ([]storage.Event, error)

type App interface {
	Create(ctx context.Context, userID int, title, desc string, start, stop time.Time, notif *time.Duration) (id int, err error)
	Update(ctx context.Context, id int, change storage.Event) error
	Delete(ctx context.Context, id int) error
	ListAll(ctx context.Context) ([]storage.Event, error)
	ListDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func New(logger logger.Logger, storage storage.Storage) App {
	return &app{
		logger,
		storage,
	}
}

var ErrNoUserID = errors.New("no user id of the event")
var ErrEmptyTitle = errors.New("no title of the event")
var ErrStartInPast = errors.New("start time of the event in the past")
var ErrDateBusy = errors.New("this time is already occupied by another event")
