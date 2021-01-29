package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type server struct {
	app    app.App
	logger logger.Logger
	srv    *http.Server
	router *mux.Router
}

func newServer(app app.App, logger logger.Logger) Server {
	s := &server{
		app:    app,
		logger: logger,
		router: mux.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *server) Start(addr string) error {
	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}

func (s *server) configureRouter() {
	s.router.Use(loggingMiddleware(s.logger))
	s.router.HandleFunc("/hello", s.handleHello).Methods(http.MethodGet)
}

func (s *server) handleHello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello, world\n"))
	if err != nil {
		s.logger.Error(fmt.Errorf("http write: %w", err))
	}
}
