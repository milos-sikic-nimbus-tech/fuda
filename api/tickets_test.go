package fuda_test

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"
)

var statusLine = regexp.MustCompile(`(?m)^\*{0,2}Status:\*{0,2}[ \t]*(.*?)[ \t]*$`)

func TestTicketStatusesAreKnownLabels(t *testing.T) {
	known := []string{"needs-triage", "needs-info", "ready-for-agent", "ready-for-human", "wontfix", "done", "claimed", "resolved"}
	tickets, err := filepath.Glob("../.scratch/*/issues/*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, ticket := range tickets {
		content, err := os.ReadFile(ticket)
		if err != nil {
			t.Fatal(err)
		}
		match := statusLine.FindSubmatch(content)
		if match == nil {
			t.Errorf("%s: no Status line", ticket)
		} else if !slices.Contains(known, string(match[1])) {
			t.Errorf("%s: unknown status %q", ticket, match[1])
		}
	}
}
