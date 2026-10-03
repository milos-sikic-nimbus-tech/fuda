package board

import (
	"slices"
	"strings"
	"unicode"

	"fuda/internal/taskfiles"
)

const (
	epicGroup  = "epic"
	plainGroup = "label"
)

var defaultTypes = []string{"bug", "feat", "impr", "refactor", "test", "chore", "docs", "question"}

func buildFacets(develop taskfiles.Result, p *people) Facets {
	return Facets{
		Statuses:    statuses(develop.Tasks),
		People:      orEmpty(p.listed),
		LabelGroups: labelGroups(develop.Config.LabelGroups, develop.Tasks),
		Prefixes:    prefixes(develop.Tasks),
	}
}

func statuses(tasks []taskfiles.Task) []string {
	var out []string
	for _, t := range tasks {
		if !slices.Contains(out, t.Status) {
			out = append(out, t.Status)
		}
	}
	return orEmpty(out)
}

func labelGroups(configured []taskfiles.LabelGroup, tasks []taskfiles.Task) []LabelGroup {
	if configured != nil {
		out := make([]LabelGroup, 0, len(configured))
		for _, g := range configured {
			out = mergeGroup(out, groupName(g.Name), g.Values)
		}
		return out
	}

	used := map[string][]string{}
	var order []string
	for _, t := range tasks {
		for _, l := range t.Labels {
			group := groupName(l.Group)
			if _, ok := used[group]; !ok {
				order = append(order, group)
			}
			if !slices.Contains(used[group], l.Value) {
				used[group] = append(used[group], l.Value)
			}
		}
	}
	if _, ok := used["type"]; !ok {
		order = append([]string{"type"}, order...)
	}
	used["type"] = mergeValues(defaultTypes, used["type"])

	out := make([]LabelGroup, 0, len(order))
	for _, group := range order {
		values := used[group]
		if group != "type" {
			slices.Sort(values)
		}
		out = append(out, LabelGroup{Name: group, Values: values})
	}
	return out
}

func mergeGroup(groups []LabelGroup, name string, values []string) []LabelGroup {
	for i := range groups {
		if groups[i].Name == name {
			groups[i].Values = mergeValues(groups[i].Values, values)
			return groups
		}
	}
	return append(groups, LabelGroup{Name: name, Values: mergeValues(nil, values)})
}

func mergeValues(base, extra []string) []string {
	out := slices.Clone(base)
	for _, v := range extra {
		if !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	return orEmpty(out)
}

func groupName(group string) string {
	switch strings.ToLower(group) {
	case "":
		return plainGroup
	case "theme", "epic":
		return epicGroup
	}
	return strings.ToLower(group)
}

func labelStrings(labels []taskfiles.Label) []string {
	out := make([]string, 0, len(labels))
	for _, l := range labels {
		out = append(out, groupName(l.Group)+":"+l.Value)
	}
	return out
}

func prefixes(tasks []taskfiles.Task) []string {
	var out []string
	for _, t := range tasks {
		if p := idPrefix(t.ID); !slices.Contains(out, p) {
			out = append(out, p)
		}
	}
	slices.Sort(out)
	return orEmpty(out)
}

func idPrefix(id string) string {
	end := strings.IndexFunc(id, func(r rune) bool { return !unicode.IsLetter(r) })
	switch {
	case end == -1:
		return id
	case end > 0:
		return id[:end]
	}
	digits := strings.IndexFunc(id, func(r rune) bool { return !unicode.IsDigit(r) })
	if digits == -1 {
		return id
	}
	return id[:digits]
}
