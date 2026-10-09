package httpapi

import (
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"fuda/internal/board"
)

type middleware func(http.Handler) http.Handler

func chain(h http.Handler, mws ...middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type GitHubLogin interface {
	Token(w http.ResponseWriter, r *http.Request) (string, error)
	Routes(mux *http.ServeMux)
}

func NewHandler(log *slog.Logger, boards *board.Boards, spa fs.FS, github GitHubLogin) http.Handler {
	mux := http.NewServeMux()
	api{log: log, boards: boards, github: github}.routes(mux)
	if github != nil {
		github.Routes(mux)
	}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/", spaHandler(spa))
	return chain(mux, recoverPanics(log), logRequests(log))
}

func logRequests(log *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			log.Info("request", "method", r.Method, "path", r.URL.Path, "status", rec.status, "duration", time.Since(start))
		})
	}
}

func recoverPanics(log *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic", "path", r.URL.Path, "error", err)
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
