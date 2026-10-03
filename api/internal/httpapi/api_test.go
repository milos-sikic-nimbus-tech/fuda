package httpapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"fuda/internal/board"
	"fuda/internal/source/local"
)

func writeRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, content := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestAPISmoke(t *testing.T) {
	root := writeRepo(t, map[string]string{
		"docs/board/tasks/SS-1-login.md": "---\nid: SS-1\ntitle: Login\nstatus: in progress\nowner: Ana\nlabels: [\"type:bug\"]\n---\nNeeds [the rules](../TASKS.md) and SS-2.\n",
		"docs/board/tasks/SS-2.md":       "---\nid: SS-2\ntitle: Logout\nstatus: backlog\n---\nRotate the session key.\n",
		"docs/board/tasks/broken.md":     "no frontmatter\n",
		"docs/board/archive/SS-0.md":     "---\nid: SS-0\ntitle: Old\nstatus: done\ndone: 2026-09-01\n---\n",
		"docs/board/TASKS.md":            "# Task rules\n\nClaim first.\n",
	})
	service := board.NewService(local.New(root, "docs", "develop"), board.Options{
		Title: "test", DocsRoot: "docs", BoardDir: "docs/board", WorkBranch: "develop", ProdBranch: "main",
		Cooldown: time.Minute,
	})
	if err := service.Sync(context.Background()); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(NewHandler(slog.New(slog.DiscardHandler), service, fstest.MapFS{}, Credentials{}, ""))
	defer server.Close()

	var b struct {
		Title    string
		Columns  []struct{ ID string }
		Cards    []struct{ ID, Column string }
		Problems []struct{ Path string }
	}
	getJSON(t, server.URL+"/api/board", http.StatusOK, &b)
	if b.Title != "test" || len(b.Columns) != 6 || len(b.Cards) != 2 || len(b.Problems) != 1 {
		t.Fatalf("board: %+v", b)
	}

	var task struct{ ID, HTML string }
	getJSON(t, server.URL+"/api/tasks/SS-1", http.StatusOK, &task)
	if !strings.Contains(task.HTML, `href="/docs/board/TASKS.md"`) {
		t.Errorf("task html: %s", task.HTML)
	}
	getJSON(t, server.URL+"/api/tasks/NOPE", http.StatusNotFound, nil)

	var ids []string
	getJSON(t, server.URL+"/api/search?q=session", http.StatusOK, &ids)
	if len(ids) != 1 || ids[0] != "SS-2" {
		t.Errorf("search: %v", ids)
	}

	var archive []struct{ ID, Done string }
	getJSON(t, server.URL+"/api/archive", http.StatusOK, &archive)
	if len(archive) != 1 || archive[0].Done != "2026-09-01" {
		t.Errorf("archive: %+v", archive)
	}

	var doc struct {
		Title     string
		Backlinks []string
	}
	getJSON(t, server.URL+"/api/docs?path=board/TASKS.md", http.StatusOK, &doc)
	if doc.Title != "Task rules" || len(doc.Backlinks) != 1 {
		t.Errorf("doc: %+v", doc)
	}
	getJSON(t, server.URL+"/api/docs?path=../../etc/passwd", http.StatusNotFound, nil)

	for _, want := range []int{http.StatusAccepted, http.StatusTooManyRequests} {
		res, err := http.Post(server.URL+"/api/sync", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != want {
			t.Errorf("sync: got %d, want %d", res.StatusCode, want)
		}
	}
}

func getJSON(t *testing.T, url string, status int, v any) {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != status {
		t.Fatalf("%s: status %d, want %d: %s", url, res.StatusCode, status, body)
	}
	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			t.Fatalf("%s: %v", url, err)
		}
	}
}

func TestBasicAuth(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	handler := basicAuth(Credentials{User: "fuda", Password: "secret"})(ok)

	cases := []struct {
		path       string
		user, pass string
		want       int
	}{
		{"/api/board", "", "", http.StatusUnauthorized},
		{"/api/board", "fuda", "wrong", http.StatusUnauthorized},
		{"/api/board", "fuda", "secret", http.StatusOK},
		{"/healthz", "", "", http.StatusOK},
		{"/api/webhooks/github", "", "", http.StatusOK},
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, c.path, nil)
		if c.user != "" {
			req.SetBasicAuth(c.user, c.pass)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != c.want {
			t.Errorf("%s as %q: got %d, want %d", c.path, c.user, rec.Code, c.want)
		}
	}
}

func TestWebhookSignature(t *testing.T) {
	sign := func(body, secret string) string {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(body))
		return "sha256=" + hex.EncodeToString(mac.Sum(nil))
	}
	cases := []struct {
		name, host, header, value string
		want                      bool
	}{
		{"github signed", "github", "X-Hub-Signature-256", sign(`{"a":1}`, "s3cret"), true},
		{"github wrong secret", "github", "X-Hub-Signature-256", sign(`{"a":1}`, "other"), false},
		{"github unsigned", "github", "", "", false},
		{"azure header", "azure", "X-Fuda-Secret", "s3cret", true},
		{"azure wrong", "azure", "X-Fuda-Secret", "nope", false},
		{"unknown host", "gitlab", "", "", false},
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodPost, "/api/webhooks/"+c.host, strings.NewReader(`{"a":1}`))
		req.SetPathValue("host", c.host)
		if c.header != "" {
			req.Header.Set(c.header, c.value)
		}
		if got := validWebhook(req, "s3cret"); got != c.want {
			t.Errorf("%s: got %v", c.name, got)
		}
	}
	if !validWebhook(httptest.NewRequest(http.MethodPost, "/api/webhooks/github", nil), "") {
		t.Error("without a secret every webhook is accepted")
	}
}
