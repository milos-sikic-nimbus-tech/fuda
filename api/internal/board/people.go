package board

import (
	"fmt"
	"slices"
	"strings"

	"fuda/internal/taskfiles"
)

type Person struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
	GitHub  string   `json:"github,omitempty"`
	Kind    string   `json:"kind"`
}

type people struct {
	configured  bool
	canonicalOf map[string]string
	byGitHub    map[string]string
	listed      []Person
}

func newPeople(configured []taskfiles.Person, tasks []taskfiles.Task) *people {
	p := &people{canonicalOf: map[string]string{}, byGitHub: map[string]string{}}
	if configured != nil {
		p.configured = true
		for _, person := range configured {
			p.listed = append(p.listed, Person{
				ID: person.ID, Name: person.Name, Aliases: orEmpty(person.Aliases), GitHub: person.GitHub, Kind: person.Kind,
			})
			p.canonicalOf[person.ID] = person.ID
			if person.GitHub != "" {
				p.byGitHub[key(person.GitHub)] = person.ID
			}
		}
		return p
	}

	var names []string
	for _, t := range tasks {
		names = append(names, t.Owners...)
		names = append(names, t.Testers...)
	}
	for _, name := range names {
		if _, seen := p.canonicalOf[key(name)]; !seen {
			p.canonicalOf[key(name)] = name
		}
	}
	p.foldFirstNames()
	var listed []string
	for _, name := range p.canonicalOf {
		if !slices.Contains(listed, name) {
			listed = append(listed, name)
		}
	}
	slices.SortFunc(listed, func(a, b string) int { return strings.Compare(key(a), key(b)) })
	for _, name := range listed {
		p.listed = append(p.listed, Person{ID: name, Name: name, Aliases: []string{}, Kind: taskfiles.KindHuman})
	}
	return p
}

func (p *people) foldFirstNames() {
	fullByFirst := map[string][]string{}
	for k, name := range p.canonicalOf {
		if first, _, multi := strings.Cut(k, " "); multi {
			fullByFirst[first] = append(fullByFirst[first], name)
		}
	}
	for k := range p.canonicalOf {
		if strings.Contains(k, " ") {
			continue
		}
		if full := fullByFirst[k]; len(full) == 1 {
			p.canonicalOf[k] = full[0]
		}
	}
}

func (p *people) canonical(names []string) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		if c, ok := p.lookup(name); ok {
			name = c
		}
		if !slices.Contains(out, name) {
			out = append(out, name)
		}
	}
	return out
}

func (p *people) lookup(name string) (string, bool) {
	if p.configured {
		id, ok := p.canonicalOf[name]
		return id, ok
	}
	c, ok := p.canonicalOf[key(name)]
	return c, ok
}

func (p *people) unknown(names []string) []string {
	if !p.configured {
		return nil
	}
	var out []string
	for _, name := range names {
		if _, ok := p.lookup(name); !ok {
			out = append(out, name)
		}
	}
	return out
}

func (p *people) problems(tasks []taskfiles.Task) []Problem {
	var out []Problem
	for _, t := range tasks {
		for _, name := range p.unknown(t.Owners) {
			out = append(out, Problem{Path: t.Path, Reason: fmt.Sprintf("owner %q is not an id in people.md", name)})
		}
		for _, name := range p.unknown(t.Testers) {
			out = append(out, Problem{Path: t.Path, Reason: fmt.Sprintf("tester %q is not an id in people.md", name)})
		}
	}
	return out
}

func (p *people) withAuthors(prs []LinkedPR) []LinkedPR {
	out := make([]LinkedPR, 0, len(prs))
	for _, pr := range prs {
		pr.AuthorID = p.byGitHub[key(pr.Author)]
		out = append(out, pr)
	}
	return out
}

func (p *people) ofGitHub(login string) string {
	return p.byGitHub[key(login)]
}

func key(name string) string {
	return strings.Join(strings.Fields(strings.ToLower(name)), " ")
}
