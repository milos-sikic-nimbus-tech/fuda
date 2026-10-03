package board

import (
	"regexp"
	"slices"
	"strings"

	"fuda/internal/taskfiles"
)

var idToken = regexp.MustCompile(`[A-Za-z0-9][A-Za-z0-9.\-]*[A-Za-z0-9]|[A-Za-z0-9]`)

type references struct {
	blockedBy   map[string][]string
	blocks      map[string][]string
	mentions    map[string][]string
	mentionedBy map[string][]string
}

func newReferences(tasks []taskfiles.Task) references {
	known := map[string]bool{}
	for _, t := range tasks {
		known[t.ID] = true
	}
	r := references{
		blockedBy:   map[string][]string{},
		blocks:      map[string][]string{},
		mentions:    map[string][]string{},
		mentionedBy: map[string][]string{},
	}
	for _, t := range tasks {
		for _, id := range KnownIDs(t.BlockedBy, known) {
			if id == t.ID {
				continue
			}
			r.blockedBy[t.ID] = append(r.blockedBy[t.ID], id)
			r.blocks[id] = append(r.blocks[id], t.ID)
		}
		for _, id := range KnownIDs(t.Body, known) {
			if id == t.ID {
				continue
			}
			r.mentions[t.ID] = append(r.mentions[t.ID], id)
			r.mentionedBy[id] = append(r.mentionedBy[id], t.ID)
		}
	}
	for _, m := range []map[string][]string{r.blockedBy, r.blocks, r.mentions, r.mentionedBy} {
		for _, t := range tasks {
			m[t.ID] = orEmpty(m[t.ID])
		}
	}
	return r
}

func KnownIDs(text string, known map[string]bool) []string {
	var out []string
	for _, token := range idToken.FindAllString(text, -1) {
		token = strings.TrimRight(token, ".-")
		if known[token] && !slices.Contains(out, token) {
			out = append(out, token)
		}
	}
	return out
}
