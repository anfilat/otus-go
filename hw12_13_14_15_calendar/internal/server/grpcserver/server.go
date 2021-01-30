package grpcserver

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type server struct {
	app    app.App
	logger logger.Logger
}

func newServer(app app.App, logger logger.Logger) *server {
	s := &server{
		app:    app,
		logger: logger,
	}
	return s
}

func (s *server) Start(addr string) error {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterCalendarServer(server, NewService())

	s.logger.Info("starting grpc server on ", addr)
	err = server.Serve(lsn)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	return nil
}
