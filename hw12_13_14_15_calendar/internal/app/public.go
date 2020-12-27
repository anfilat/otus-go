package app

import (
	"context"
	"errors"
	"time"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App interface {
	CreateEvent(ctx context.Context, userID int, title, desc string, start, stop time.Time, notif *time.Duration) (id int, err error)
	UpdateEvent(ctx context.Context, id int, change storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
	ListAllEvents(ctx context.Context) ([]storage.Event, error)
	ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
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
