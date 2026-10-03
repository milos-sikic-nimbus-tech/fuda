package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"fuda/internal/board"
	"fuda/internal/guide"
)

type api struct {
	log           *slog.Logger
	service       *board.Service
	webhookSecret string
}

func (a api) routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/board", a.board)
	mux.HandleFunc("GET /api/tasks/{id}", a.task)
	mux.HandleFunc("GET /api/search", a.search)
	mux.HandleFunc("GET /api/archive", a.archive)
	mux.HandleFunc("GET /api/docs", a.doc)
	mux.HandleFunc("GET /api/guide", a.guideList)
	mux.HandleFunc("GET /api/guide/{slug}", a.guidePage)
	mux.HandleFunc("POST /api/sync", a.sync)
	mux.HandleFunc("POST /api/webhooks/{host}", a.webhook)
}

func (a api) board(w http.ResponseWriter, _ *http.Request) {
	view, ready := a.service.Board()
	if !ready {
		a.json(w, http.StatusServiceUnavailable, view)
		return
	}
	a.json(w, http.StatusOK, view)
}

func (a api) task(w http.ResponseWriter, r *http.Request) {
	view, err := a.service.Task(r.PathValue("id"))
	a.result(w, view, err)
}

func (a api) search(w http.ResponseWriter, r *http.Request) {
	a.json(w, http.StatusOK, a.service.Search(r.URL.Query().Get("q")))
}

func (a api) archive(w http.ResponseWriter, _ *http.Request) {
	a.json(w, http.StatusOK, a.service.Archive())
}

func (a api) doc(w http.ResponseWriter, r *http.Request) {
	view, err := a.service.Doc(r.URL.Query().Get("path"))
	a.result(w, view, err)
}

func (a api) guideList(w http.ResponseWriter, _ *http.Request) {
	a.json(w, http.StatusOK, guide.List())
}

func (a api) guidePage(w http.ResponseWriter, r *http.Request) {
	page, err := guide.Render(r.PathValue("slug"))
	if errors.Is(err, guide.ErrNotFound) {
		err = board.ErrNotFound
	}
	a.result(w, page, err)
}

func (a api) sync(w http.ResponseWriter, r *http.Request) {
	if !a.service.RequestSync(r.Context()) {
		a.json(w, http.StatusTooManyRequests, map[string]string{"error": "a sync ran moments ago; try again shortly"})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (a api) webhook(w http.ResponseWriter, r *http.Request) {
	if !validWebhook(r, a.webhookSecret) {
		a.json(w, http.StatusUnauthorized, map[string]string{"error": "invalid webhook signature"})
		return
	}
	a.service.NotifyChange(r.Context())
	w.WriteHeader(http.StatusAccepted)
}

func (a api) result(w http.ResponseWriter, v any, err error) {
	switch {
	case errors.Is(err, board.ErrNotFound):
		a.json(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case err != nil:
		a.log.Error("request failed", "error", err)
		a.json(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	default:
		a.json(w, http.StatusOK, v)
	}
}

func (a api) json(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		a.log.Error("encode response", "error", err)
	}
}
