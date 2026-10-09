package board

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

type fakeRepos map[string]map[string][]byte

type fakeSource struct{ files map[string][]byte }

func (f fakeSource) Head(context.Context, string) (string, error) { return "head", nil }

func (f fakeSource) Files(context.Context, string) (map[string][]byte, error) { return f.files, nil }

type missingSource struct{}

func (missingSource) Head(context.Context, string) (string, error) {
	return "", ErrBranchMissing
}

func (missingSource) Files(context.Context, string) (map[string][]byte, error) {
	return nil, ErrBranchMissing
}

func taskFile(id string) map[string][]byte {
	return map[string][]byte{"docs/board/tasks/" + id + ".md": []byte("---\nid: " + id + "\ntitle: " + id + "\nstatus: backlog\n---\n")}
}

func newTestBoards(t *testing.T, repos fakeRepos) *Boards {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return NewBoards(ctx, slog.New(slog.DiscardHandler), time.Hour, func(id BoardID) (*Service, error) {
		files, ok := repos[id.Repo]
		if !ok {
			return nil, ErrNotFound
		}
		var source Source = fakeSource{files}
		if files == nil {
			source = missingSource{}
		}
		return NewService(source, Options{DocsRoot: "docs", BoardDir: "docs/board", WorkBranch: "develop", ProdBranch: "main", Origin: Origin{Path: id.Path()}}), nil
	})
}

func TestTwoBoardsKeepTheirOwnTasks(t *testing.T) {
	boards := newTestBoards(t, fakeRepos{"a/one": taskFile("ONE-1"), "a/two": taskFile("TWO-1")})

	one, err := boards.Get(BoardID{"github", "a/one"})
	if err != nil {
		t.Fatal(err)
	}
	two, err := boards.Get(BoardID{"github", "a/two"})
	if err != nil {
		t.Fatal(err)
	}

	viewOne, _ := one.Board()
	viewTwo, _ := two.Board()
	if len(viewOne.Cards) != 1 || viewOne.Cards[0].ID != "ONE-1" {
		t.Errorf("first board: %+v", viewOne.Cards)
	}
	if len(viewTwo.Cards) != 1 || viewTwo.Cards[0].ID != "TWO-1" {
		t.Errorf("second board: %+v", viewTwo.Cards)
	}
	if _, err := two.Task("ONE-1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("a task of the first board is visible in the second: %v", err)
	}
}

func TestGetReturnsTheSameServiceEachTime(t *testing.T) {
	boards := newTestBoards(t, fakeRepos{"a/one": taskFile("ONE-1")})
	first, _ := boards.Get(BoardID{"github", "a/one"})
	second, _ := boards.Get(BoardID{"github", "a/one"})
	if first != second {
		t.Error("a Board was opened twice")
	}
}

func TestGetUnknownBoardIsNotFound(t *testing.T) {
	boards := newTestBoards(t, fakeRepos{"a/gone": nil})
	for _, id := range []BoardID{
		{"github", "a/nothing"},
		{"github", "a/gone"},
		{"github", "../etc"},
		{"github", ""},
	} {
		if _, err := boards.Get(id); !errors.Is(err, ErrNotFound) {
			t.Errorf("%v: got %v, want not found", id, err)
		}
	}
}
