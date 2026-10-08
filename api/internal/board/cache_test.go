package board

import (
	"context"
	"errors"
	"testing"
)

type fixedSource struct {
	files map[string][]byte
	err   error
}

func (f fixedSource) Head(context.Context, string) (string, error) { return "sha1", f.err }

func (f fixedSource) Files(context.Context, string) (map[string][]byte, error) { return f.files, f.err }

func TestRestoreServesTheLastSyncedBoard(t *testing.T) {
	dir := t.TempDir()
	opts := Options{WorkBranch: "develop", BoardDir: "docs/board", DocsRoot: "docs", CacheDir: dir}
	files := map[string][]byte{
		"docs/board/tasks/A1.md": []byte("---\nid: A1\ntitle: First\nstatus: backlog\n---\nBody\n"),
	}
	if err := NewService(fixedSource{files: files}, opts).Sync(context.Background()); err != nil {
		t.Fatal(err)
	}

	offline := NewService(fixedSource{err: errors.New("host is down")}, opts)
	if err := offline.Restore(); err != nil {
		t.Fatal(err)
	}
	view, ready := offline.Board()
	if !ready || len(view.Cards) != 1 || view.Cards[0].ID != "A1" {
		t.Fatalf("restored board: ready=%v cards=%+v", ready, view.Cards)
	}
	if view.Sync.Develop.SHA != "sha1" {
		t.Errorf("restored status keeps the cached head, got %q", view.Sync.Develop.SHA)
	}
	if err := offline.Sync(context.Background()); err == nil {
		t.Error("a failing host still reports the error")
	}
	if _, ready := offline.Board(); !ready {
		t.Error("a failed sync keeps the restored board")
	}
}

func TestRestoreWithoutCacheStartsEmpty(t *testing.T) {
	s := NewService(fixedSource{}, Options{CacheDir: t.TempDir()})
	if err := s.Restore(); err != nil {
		t.Fatal(err)
	}
	if _, ready := s.Board(); ready {
		t.Error("nothing cached means no board yet")
	}
}
