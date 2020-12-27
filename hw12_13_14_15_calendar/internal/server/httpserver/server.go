package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
		Addr:    addr,
		Handler: s.router,
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

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.loggingMiddleware)
	s.router.HandleFunc("/hello", s.handleHello).Methods("GET")
}

func (s *server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		s.logger.Info(
			fmt.Sprintf("%s %s %s %s %d %s %s",
				requestAddr(r),
				r.Method,
				r.RequestURI,
				r.Proto,
				rw.code,
				latency(start),
				userAgent(r),
			))
	})
}

func requestAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func userAgent(r *http.Request) string {
	userAgents := r.Header["User-Agent"]
	if len(userAgents) > 0 {
		return "\"" + userAgents[0] + "\""
	}
	return ""
}

func latency(start time.Time) string {
	return fmt.Sprintf("%dms", time.Since(start).Milliseconds())
}

func (s *server) handleHello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello, world\n"))
	if err != nil {
		s.logger.Error(fmt.Errorf("http write: %w", err))
	}
}
