package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	Base
	Events
}

type Base interface {
	Connect(ctx context.Context, connect string) error
	Close(ctx context.Context) error
}

type Events interface {
	Create(ctx context.Context, event Event) (int, error)
	Update(ctx context.Context, id int, change Event) error
	Delete(ctx context.Context, id int) error
	DeleteAll(ctx context.Context) error
	ListAll(ctx context.Context) ([]Event, error)
	ListDay(ctx context.Context, date time.Time) ([]Event, error)
	ListWeek(ctx context.Context, date time.Time) ([]Event, error)
	ListMonth(ctx context.Context, date time.Time) ([]Event, error)
	IsTimeBusy(ctx context.Context, start, stop time.Time, excludeID int) (bool, error)
}

type Event struct {
	ID           int
	Title        string
	Start        time.Time
	Stop         time.Time
	Description  string
	UserID       int
	Notification *time.Duration
}

var ErrNotExistsEvent = errors.New("no such event")
