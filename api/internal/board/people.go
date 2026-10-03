package board

import (
	"slices"
	"strings"

	"fuda/internal/taskfiles"
)

type people struct {
	canonicalOf map[string]string
	listed      []string
}

func newPeople(configured []taskfiles.Person, tasks []taskfiles.Task) *people {
	p := &people{canonicalOf: map[string]string{}}
	if configured != nil {
		for _, person := range configured {
			p.listed = append(p.listed, person.Name)
			p.canonicalOf[key(person.Name)] = person.Name
			for _, alias := range person.Aliases {
				p.canonicalOf[key(alias)] = person.Name
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
	for _, name := range p.canonicalOf {
		if !slices.Contains(p.listed, name) {
			p.listed = append(p.listed, name)
		}
	}
	slices.SortFunc(p.listed, func(a, b string) int { return strings.Compare(key(a), key(b)) })
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
		if c, ok := p.canonicalOf[key(name)]; ok {
			name = c
		}
		if !slices.Contains(out, name) {
			out = append(out, name)
		}
	}
	return out
}

func key(name string) string {
	return strings.Join(strings.Fields(strings.ToLower(name)), " ")
}
