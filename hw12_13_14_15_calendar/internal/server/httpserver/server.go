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
	router := s.router
	router.Use(loggingMiddleware(s.logger))

	router.HandleFunc("/hello", s.handleHello).Methods(http.MethodGet)

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/create", s.handleCreate).Methods(http.MethodPost)
	apiRouter.HandleFunc("/update", s.handleUpdate).Methods(http.MethodPost)
	apiRouter.HandleFunc("/delete", s.handleDelete).Methods(http.MethodPost)
	apiRouter.HandleFunc("/listday", s.handleListDay).Methods(http.MethodPost)
	apiRouter.HandleFunc("/listweek", s.handleListWeek).Methods(http.MethodPost)
	apiRouter.HandleFunc("/listmonth", s.handleListMonth).Methods(http.MethodPost)
}

func (s *server) handleHello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, []byte("Hello, world\n"))
}

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

	change := httpEventToStorageEvent(req)
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
		result = append(result, storageEventToHTTPEvent(event))
	}
	s.writeJSON(w, result)
}

func (s *server) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	data, _ := json.Marshal(v)
	fmt.Fprint(w, data)
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
