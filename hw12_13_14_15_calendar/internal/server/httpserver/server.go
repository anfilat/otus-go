package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type server struct {
	app    app.App
	logger logger.Logger
	srv    *http.Server
	router *mux.Router
}

func newServer(app app.App, logger logger.Logger) *server {
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
	router := s.router
	router.Use(loggingMiddleware(s.logger))

	router.HandleFunc("/hello", handleHello).Methods(http.MethodGet)

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/create", handleCreate(s.app)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/update", handleUpdate(s.app)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/delete", handleDelete(s.app)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/listday", handleListDay(s.app)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/listweek", handleListWeek(s.app)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/listmonth", handleListMonth(s.app)).Methods(http.MethodPost)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	data, _ := json.Marshal(v)
	//nolint:errcheck
	w.Write(data)
}

func httpEventToStorageEvent(event Event) storage.Event {
	return storage.Event{
		ID:           event.ID,
		Title:        event.Title,
		Start:        event.Start,
		Stop:         event.Stop,
		Description:  event.Description,
		UserID:       event.UserID,
		Notification: event.Notification,
	}
}

func storageEventToHTTPEvent(event storage.Event) Event {
	return Event{
		ID:           event.ID,
		Title:        event.Title,
		Start:        event.Start,
		Stop:         event.Stop,
		Description:  event.Description,
		UserID:       event.UserID,
		Notification: event.Notification,
	}
}
