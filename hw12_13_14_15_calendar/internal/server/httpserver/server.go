package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	s.router.Use(loggingMiddleware(s.logger))
	s.router.HandleFunc("/hello", s.handleHello).Methods(http.MethodGet)
	s.router.HandleFunc("/api/create", s.handleCreate).Methods(http.MethodPost)
	s.router.HandleFunc("/api/update", s.handleUpdate).Methods(http.MethodPost)
	s.router.HandleFunc("/api/delete", s.handleDelete).Methods(http.MethodPost)
	s.router.HandleFunc("/api/listday", s.handleListDay).Methods(http.MethodPost)
	s.router.HandleFunc("/api/listweek", s.handleListWeek).Methods(http.MethodPost)
	s.router.HandleFunc("/api/listmonth", s.handleListMonth).Methods(http.MethodPost)
}

func (s *server) handleHello(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello, world\n"))
	if err != nil {
		s.logger.Error(fmt.Errorf("http write: %w", err))
	}
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

func (s *server) handleCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := Event{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.app.Create(r.Context(), req.UserID, req.Title, req.Description, req.Start, req.Stop, req.Notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.writeJSON(w, CreateResult{id})
}

func (s *server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := Event{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	change := storage.Event{
		ID:           req.ID,
		Title:        req.Title,
		Start:        req.Start,
		Stop:         req.Stop,
		Description:  req.Description,
		UserID:       req.UserID,
		Notification: req.Notification,
	}
	err = s.app.Update(r.Context(), req.ID, change)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.writeJSON(w, OkResult{Ok: true})
}

func (s *server) handleDelete(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := DeleteRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Delete(r.Context(), req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.writeJSON(w, OkResult{Ok: true})
}

func (s *server) handleListDay(w http.ResponseWriter, r *http.Request) {
	s.handleList(w, r, s.app.ListDay)
}

func (s *server) handleListWeek(w http.ResponseWriter, r *http.Request) {
	s.handleList(w, r, s.app.ListMonth)
}

func (s *server) handleListMonth(w http.ResponseWriter, r *http.Request) {
	s.handleList(w, r, s.app.ListWeek)
}

func (s *server) handleList(w http.ResponseWriter, r *http.Request, fn app.ListEvents) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := ListRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := fn(r.Context(), req.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := make(ListResult, 0, len(events))
	for _, event := range events {
		result = append(result, Event{
			ID:           event.ID,
			Title:        event.Title,
			Start:        event.Start,
			Stop:         event.Stop,
			Description:  event.Description,
			UserID:       event.UserID,
			Notification: event.Notification,
		})
	}
	s.writeJSON(w, result)
}

func (s *server) writeJSON(w http.ResponseWriter, v interface{}) {
	data, _ := json.Marshal(v)
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		s.logger.Error(fmt.Errorf("http write: %w", err))
	}
}
