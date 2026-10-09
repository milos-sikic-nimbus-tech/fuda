package taskfiles

import (
	"errors"
	"testing"
)

func TestApplyMove(t *testing.T) {
	tests := []struct {
		name string
		in   string
		move Move
		want string
	}{
		{
			name: "changes only the status line",
			in:   "---\nid: T-1\ntitle: One\nstatus: backlog\nlabels: a, b\n---\n\nBody: keep\n",
			move: Move{Status: "in progress"},
			want: "---\nid: T-1\ntitle: One\nstatus: in progress\nlabels: a, b\n---\n\nBody: keep\n",
		},
		{
			name: "adds claimed right after status",
			in:   "---\nid: T-1\ntitle: One\nstatus: backlog\nadded: 2026-01-01\n---\nBody\n",
			move: Move{Status: "in progress", Claimed: "2026-10-09"},
			want: "---\nid: T-1\ntitle: One\nstatus: in progress\nclaimed: 2026-10-09\nadded: 2026-01-01\n---\nBody\n",
		},
		{
			name: "never changes a claimed date that is set",
			in:   "---\nid: T-1\ntitle: One\nstatus: testing\nclaimed: 2026-02-02\n---\nBody\n",
			move: Move{Status: "in progress", Claimed: "2026-10-09"},
			want: "---\nid: T-1\ntitle: One\nstatus: in progress\nclaimed: 2026-02-02\n---\nBody\n",
		},
		{
			name: "fills an empty claimed line",
			in:   "---\nid: T-1\ntitle: One\nstatus: backlog\nclaimed:\n---\n",
			move: Move{Status: "in progress", Claimed: "2026-10-09"},
			want: "---\nid: T-1\ntitle: One\nstatus: in progress\nclaimed: 2026-10-09\n---\n",
		},
		{
			name: "keeps windows line endings and the byte order mark",
			in:   "\xef\xbb\xbf---\r\nid: T-1\r\ntitle: One\r\nstatus: backlog\r\n---\r\nBody\r\n",
			move: Move{Status: "in progress", Claimed: "2026-10-09"},
			want: "\xef\xbb\xbf---\r\nid: T-1\r\ntitle: One\r\nstatus: in progress\r\nclaimed: 2026-10-09\r\n---\r\nBody\r\n",
		},
		{
			name: "quotes a status that needs it",
			in:   "---\nid: T-1\ntitle: One\nstatus: backlog\n---\n",
			move: Move{Status: "in: review"},
			want: "---\nid: T-1\ntitle: One\nstatus: 'in: review'\n---\n",
		},
		{
			name: "ignores status lines in the body and nested keys",
			in:   "---\nid: T-1\ntitle: One\nmeta:\n  status: nested\nstatus: backlog\n---\nstatus: body\n",
			move: Move{Status: "testing"},
			want: "---\nid: T-1\ntitle: One\nmeta:\n  status: nested\nstatus: testing\n---\nstatus: body\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ApplyMove([]byte(tt.in), tt.move)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestApplyMoveNeedsAStatusLine(t *testing.T) {
	for name, in := range map[string]string{
		"no frontmatter":      "Body only\n",
		"no status line":      "---\nid: T-1\ntitle: One\n---\n",
		"status only in body": "---\nid: T-1\n---\nstatus: backlog\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ApplyMove([]byte(in), Move{Status: "testing"}); !errors.Is(err, ErrNoStatusLine) {
				t.Errorf("got %v, want ErrNoStatusLine", err)
			}
		})
	}
}

func TestParseOne(t *testing.T) {
	task, err := ParseOne("docs/board/tasks/T-1.md", []byte("---\nid: T-1\ntitle: One\nstatus: backlog\n---\n"))
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != "backlog" || task.ID != "T-1" {
		t.Errorf("got %+v", task)
	}
	if _, err := ParseOne("docs/board/tasks/T-1.md", []byte("nope")); err == nil {
		t.Error("want an error for a file without frontmatter")
	}
}
