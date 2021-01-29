package httpserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(rw, r)

			logger.Info(
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
