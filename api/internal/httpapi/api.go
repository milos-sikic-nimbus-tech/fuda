package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"fuda/internal/board"
	"fuda/internal/guide"
)

type api struct {
	log           *slog.Logger
	boards        *board.Boards
	webhookSecret string
}

type boardHost struct {
	name   string
	params []string
}

var boardHosts = []boardHost{
	{"github", []string{"owner", "repo"}},
	{"azure", []string{"org", "project", "repo"}},
	{"local", []string{"folder"}},
}

type boardHandler func(w http.ResponseWriter, r *http.Request, service *board.Service)

func (a api) routes(mux *http.ServeMux) {
	for _, host := range boardHosts {
		prefix := "/api/" + host.name
		for _, param := range host.params {
			prefix += "/{" + param + "}"
		}
		for pattern, handler := range map[string]boardHandler{
			"GET /board":      a.board,
			"GET /tasks/{id}": a.task,
			"GET /search":     a.search,
			"GET /archive":    a.archive,
			"GET /docs":       a.doc,
			"GET /files":      a.file,
			"POST /sync":      a.sync,
		} {
			method, path, _ := strings.Cut(pattern, " ")
			mux.HandleFunc(method+" "+prefix+path, a.onBoard(host, handler))
		}
	}
	mux.HandleFunc("GET /api/guide", a.guideList)
	mux.HandleFunc("GET /api/guide/{slug}", a.guidePage)
	mux.HandleFunc("POST /api/webhooks/{host}", a.webhook)
}

func (a api) onBoard(host boardHost, handler boardHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := make([]string, len(host.params))
		for i, param := range host.params {
			parts[i] = r.PathValue(param)
		}
		service, err := a.boards.Get(board.BoardID{Host: host.name, Repo: strings.Join(parts, "/")})
		if err != nil {
			a.result(w, nil, err)
			return
		}
		handler(w, r, service)
	}
}

func (a api) board(w http.ResponseWriter, _ *http.Request, service *board.Service) {
	view, ready := service.Board()
	if !ready {
		a.json(w, http.StatusServiceUnavailable, view)
		return
	}
	a.json(w, http.StatusOK, view)
}

func (a api) task(w http.ResponseWriter, r *http.Request, service *board.Service) {
	view, err := service.Task(r.PathValue("id"))
	a.result(w, view, err)
}

func (a api) search(w http.ResponseWriter, r *http.Request, service *board.Service) {
	a.json(w, http.StatusOK, service.Search(r.URL.Query().Get("q")))
}

func (a api) archive(w http.ResponseWriter, _ *http.Request, service *board.Service) {
	a.json(w, http.StatusOK, service.Archive())
}

func (a api) doc(w http.ResponseWriter, r *http.Request, service *board.Service) {
	view, err := service.Doc(r.URL.Query().Get("path"))
	a.result(w, view, err)
}

func (a api) file(w http.ResponseWriter, r *http.Request, service *board.Service) {
	content, contentType, err := service.Asset(r.URL.Query().Get("path"))
	if err != nil {
		a.result(w, nil, err)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")
	w.Header().Set("Cache-Control", "private, max-age=300")
	_, _ = w.Write(content)
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

func (a api) sync(w http.ResponseWriter, r *http.Request, service *board.Service) {
	if !service.RequestSync(r.Context()) {
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
	a.boards.NotifyHost(r.Context(), r.PathValue("host"))
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
