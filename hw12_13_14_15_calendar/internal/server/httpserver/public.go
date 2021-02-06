package httpserver

import (
	"context"
	"time"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type Server interface {
	Start(addr string) error
	Stop(ctx context.Context) error
}

func NewServer(app app.App, logger logger.Logger) Server {
	return newServer(app, logger)
}

type Event struct {
	ID           int
	Title        string
	Start        time.Time
	Stop         time.Time
	Description  string
	UserID       int
	Notification *time.Duration `json:"notification,omitempty"`
}

type DeleteRequest struct {
	ID int
}

type ListRequest struct {
	Date time.Time
}

type CreateResult struct {
	ID int
}

type OkResult struct {
	Ok bool
}

type ListResult []Event
