package board

import (
	"context"
	"testing"
)

func TestMeMapsTheLoginToAPerson(t *testing.T) {
	host := newFakeHost(map[string]string{
		taskPath:               taskContent("backlog", ""),
		"docs/board/people.md": "---\npeople:\n  - id: ann\n    name: Ann\n    github: ann-gh\n---\n",
	})
	host.login = "Ann-GH"
	s := newMoveService(t, host)

	me, err := s.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	check(t, "me", *me, Me{Login: "Ann-GH", PersonID: "ann"})

	host.login = "someone-else"
	me, _ = s.Me(context.Background())
	check(t, "cached for the token", me.Login, "Ann-GH")
}

func TestMeIsEmptyWithoutAViewerLogin(t *testing.T) {
	s := NewService(fakeLocal{}, Options{BoardDir: "docs/board"})

	me, err := s.Me(context.Background())
	if err != nil || me != nil {
		t.Fatalf("got %v, %v; want nil, nil", me, err)
	}
}

type fakeLocal struct{}

func (fakeLocal) Head(context.Context, string) (string, error) { return "1", nil }

func (fakeLocal) Files(context.Context, string) (map[string][]byte, error) {
	return map[string][]byte{}, nil
}
