package board

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type countingSource struct{ heads atomic.Int32 }

func (c *countingSource) Head(context.Context, string) (string, error) {
	c.heads.Add(1)
	return "sha", nil
}

func (c *countingSource) Files(context.Context, string) (map[string][]byte, error) {
	return map[string][]byte{}, nil
}

func TestNotifyChangeDefersInsteadOfDropping(t *testing.T) {
	source := &countingSource{}
	s := NewService(source, Options{WorkBranch: "develop", BoardDir: "docs/board", DocsRoot: "docs", Cooldown: 80 * time.Millisecond})

	if !s.RequestSync(context.Background()) {
		t.Fatal("first manual sync must be accepted")
	}
	s.NotifyChange(context.Background())
	s.NotifyChange(context.Background())

	waitFor(t, func() bool { return source.heads.Load() == 1 })
	time.Sleep(20 * time.Millisecond)
	if got := source.heads.Load(); got != 1 {
		t.Fatalf("webhooks inside the cooldown must wait, got %d syncs", got)
	}
	waitFor(t, func() bool { return source.heads.Load() == 2 })
	time.Sleep(120 * time.Millisecond)
	if got := source.heads.Load(); got != 2 {
		t.Fatalf("two webhooks in one cooldown collapse into one deferred sync, got %d", got)
	}
}

func waitFor(t *testing.T, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !ok() {
		if time.Now().After(deadline) {
			t.Fatal("timed out")
		}
		time.Sleep(5 * time.Millisecond)
	}
}
