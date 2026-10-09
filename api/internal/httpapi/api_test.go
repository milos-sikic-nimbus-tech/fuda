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
		"docs/board/assets/flow.png":     "\x89PNG fake",
	})
	server := httptest.NewServer(NewHandler(slog.New(slog.DiscardHandler), testBoards(t, map[string]string{"test": root}), "/local/test", fstest.MapFS{}, Credentials{}, ""))
	defer server.Close()
	base := server.URL + "/api/local/test"

	var b struct {
		Title    string
		Columns  []struct{ ID string }
		Cards    []struct{ ID, Column string }
		Problems []struct{ Path string }
	}
	getJSON(t, base+"/board", http.StatusOK, &b)
	if b.Title != "test" || len(b.Columns) != 6 || len(b.Cards) != 2 || len(b.Problems) != 1 {
		t.Fatalf("board: %+v", b)
	}

	var task struct{ ID, HTML string }
	getJSON(t, base+"/tasks/SS-1", http.StatusOK, &task)
	if !strings.Contains(task.HTML, `href="/local/test/docs/board/TASKS.md"`) {
		t.Errorf("task html: %s", task.HTML)
	}
	getJSON(t, base+"/tasks/NOPE", http.StatusNotFound, nil)

	var ids []string
	getJSON(t, base+"/search?q=session", http.StatusOK, &ids)
	if len(ids) != 1 || ids[0] != "SS-2" {
		t.Errorf("search: %v", ids)
	}

	var archive []struct{ ID, Done string }
	getJSON(t, base+"/archive", http.StatusOK, &archive)
	if len(archive) != 1 || archive[0].Done != "2026-09-01" {
		t.Errorf("archive: %+v", archive)
	}

	var doc struct {
		Title     string
		Backlinks []string
	}
	getJSON(t, base+"/docs?path=board/TASKS.md", http.StatusOK, &doc)
	if doc.Title != "Task rules" || len(doc.Backlinks) != 1 {
		t.Errorf("doc: %+v", doc)
	}
	getJSON(t, base+"/docs?path=../../etc/passwd", http.StatusNotFound, nil)

	res, err := http.Get(base + "/files?path=docs/board/assets/flow.png")
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "image/png" {
		t.Errorf("asset: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	getJSON(t, base+"/files?path=docs/board/TASKS.md", http.StatusNotFound, nil)

	for _, want := range []int{http.StatusAccepted, http.StatusTooManyRequests} {
		res, err := http.Post(base+"/sync", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != want {
			t.Errorf("sync: got %d, want %d", res.StatusCode, want)
		}
	}
}

func testBoards(t *testing.T, folders map[string]string) *board.Boards {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return board.NewBoards(ctx, slog.New(slog.DiscardHandler), time.Hour, func(id board.BoardID) (*board.Service, error) {
		root, ok := folders[id.Repo]
		if !ok || id.Host != "local" {
			return nil, board.ErrNotFound
		}
		return board.NewService(local.New(root, "docs", "develop"), board.Options{
			Title: id.Repo, DocsRoot: "docs", BoardDir: "docs/board", WorkBranch: "develop", ProdBranch: "main",
			Cooldown: time.Minute, Origin: board.Origin{Host: "local", Repo: id.Repo, Path: id.Path()},
		}), nil
	})
}

func TestBoardsAreServedByPath(t *testing.T) {
	one := writeRepo(t, map[string]string{"docs/board/tasks/ONE-1.md": "---\nid: ONE-1\ntitle: First\nstatus: backlog\n---\nSee [rules](../TASKS.md).\n", "docs/board/TASKS.md": "# Rules\n"})
	two := writeRepo(t, map[string]string{"docs/board/tasks/TWO-1.md": "---\nid: TWO-1\ntitle: Second\nstatus: backlog\n---\n"})
	server := httptest.NewServer(NewHandler(slog.New(slog.DiscardHandler), testBoards(t, map[string]string{"one": one, "two": two}), "", fstest.MapFS{}, Credentials{}, ""))
	defer server.Close()

	for folder, card := range map[string]string{"one": "ONE-1", "two": "TWO-1"} {
		var b struct{ Cards []struct{ ID string } }
		getJSON(t, server.URL+"/api/local/"+folder+"/board", http.StatusOK, &b)
		if len(b.Cards) != 1 || b.Cards[0].ID != card {
			t.Errorf("%s: %+v", folder, b.Cards)
		}
	}
	getJSON(t, server.URL+"/api/local/two/tasks/ONE-1", http.StatusNotFound, nil)

	var task struct{ HTML string }
	getJSON(t, server.URL+"/api/local/one/tasks/ONE-1", http.StatusOK, &task)
	if !strings.Contains(task.HTML, `href="/local/one/docs/board/TASKS.md"`) {
		t.Errorf("links are not prefixed with the Board path: %s", task.HTML)
	}

	getJSON(t, server.URL+"/api/local/nothing/board", http.StatusNotFound, nil)
	getJSON(t, server.URL+"/api/github/nobody/none/board", http.StatusNotFound, nil)
}

func TestRootRedirectsToTheDefaultBoard(t *testing.T) {
	server := httptest.NewServer(NewHandler(slog.New(slog.DiscardHandler), testBoards(t, nil), "/local/test", fstest.MapFS{}, Credentials{}, ""))
	defer server.Close()
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != "/local/test" {
		t.Errorf("got %d to %q", res.StatusCode, res.Header.Get("Location"))
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
