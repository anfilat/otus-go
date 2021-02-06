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
	srv    *grpc.Server
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

	s.srv = grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor(s.logger)))
	RegisterCalendarServer(s.srv, NewService(s.app))

	s.logger.Info("starting grpc server on ", addr)
	return s.srv.Serve(lsn)
}

func (s *server) Stop(_ context.Context) error {
	s.srv.GracefulStop()
	return nil
}
