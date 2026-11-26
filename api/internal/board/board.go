package board

import (
	"cmp"
	"slices"
	"time"

	"fuda/internal/taskfiles"
)

type Board struct {
	Columns  []Column  `json:"columns"`
	Cards    []Card    `json:"cards"`
	Facets   Facets    `json:"facets"`
	Maps     []EpicMap `json:"maps"`
	Problems []Problem `json:"problems"`
}

type Column struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Statuses []string `json:"statuses"`
	PROpen   bool     `json:"prOpen"`
	Done     bool     `json:"done"`
	Days     int      `json:"days,omitempty"`
	Unknown  bool     `json:"unknown"`
	Activity bool     `json:"activity"`
}

type Card struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Status       string              `json:"status"`
	Column       string              `json:"column"`
	Owners       []string            `json:"owners"`
	Testers      []string            `json:"testers"`
	Labels       []string            `json:"labels"`
	Prefix       string              `json:"prefix"`
	Added        string              `json:"added"`
	Claimed      string              `json:"claimed"`
	Done         string              `json:"done,omitempty"`
	PRs          []int               `json:"prs"`
	PRRef        string              `json:"prRef,omitempty"`
	OpenPRs      []LinkedPR          `json:"openPrs"`
	BlockedBy    string              `json:"blockedBy,omitempty"`
	BlockedByIDs []string            `json:"blockedByIds"`
	Blocks       []string            `json:"blocks"`
	References   []string            `json:"references"`
	ReferencedBy []string            `json:"referencedBy"`
	Activity     *taskfiles.Activity `json:"activity,omitempty"`
	Archived     bool                `json:"archived"`
	Path         string              `json:"path"`
}

type EpicMap struct {
	Epic     string         `json:"epic"`
	Body     string         `json:"body"`
	Counts   map[string]int `json:"counts"`
	Archived int            `json:"archived"`
}

type Facets struct {
	Statuses    []string          `json:"statuses"`
	People      []Person          `json:"people"`
	Testers     []string          `json:"testers"`
	LabelGroups []LabelGroup      `json:"labelGroups"`
	LabelColors map[string]string `json:"labelColors"`
	Prefixes    []string          `json:"prefixes"`
}

type LabelGroup struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Problem struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type Inputs struct {
	Parsed  taskfiles.Result
	OpenPRs map[string][]LinkedPR
	Today   time.Time
}

func Build(in Inputs) Board {
	tasks := in.Parsed.Tasks
	columns := buildColumns(in.Parsed.Config.Stages, tasks)
	done := recentlyDone(columns.done, in.Parsed.Archived, in.Today)
	shown := slices.Concat(tasks, done)
	people := newPeople(in.Parsed.Config.People, shown)
	refs := newReferences(tasks)

	cards := make([]Card, 0, len(shown))
	for _, t := range tasks {
		openPRs := in.OpenPRs[t.ID]
		column := columns.place(t.Status, len(openPRs) > 0)
		cards = append(cards, Card{
			ID:           t.ID,
			Title:        t.Title,
			Status:       t.Status,
			Column:       column,
			Owners:       people.canonical(t.Owners),
			Testers:      people.canonical(t.Testers),
			Labels:       labelStrings(t.Labels),
			Prefix:       idPrefix(t.ID),
			Added:        t.Added,
			Claimed:      t.Claimed,
			PRs:          orEmpty(t.PRs),
			PRRef:        t.PRRef,
			OpenPRs:      people.withAuthors(openPRs),
			BlockedBy:    t.BlockedBy,
			BlockedByIDs: refs.blockedBy[t.ID],
			Blocks:       refs.blocks[t.ID],
			References:   refs.mentions[t.ID],
			ReferencedBy: refs.mentionedBy[t.ID],
			Activity:     columns.latestActivity(column, t.Activities),
			Path:         t.Path,
		})
	}
	for _, t := range done {
		c := archivedCard(t)
		c.Column = columns.done.ID
		c.Owners = people.canonical(t.Owners)
		c.Testers = people.canonical(t.Testers)
		cards = append(cards, c)
	}

	facetsFrom := in.Parsed
	facetsFrom.Tasks = shown
	return Board{
		Columns:  columns.list,
		Cards:    cards,
		Facets:   buildFacets(facetsFrom, people),
		Maps:     buildMaps(in.Parsed, cards),
		Problems: append(problems(in.Parsed.Problems), people.problems(tasks)...),
	}
}

func recentlyDone(column *Column, archived []taskfiles.Task, today time.Time) []taskfiles.Task {
	if column == nil {
		return nil
	}
	day := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	first := day.AddDate(0, 0, 1-column.Days)
	var out []taskfiles.Task
	for _, t := range archived {
		done, err := time.Parse(time.DateOnly, t.Done)
		if err != nil || done.Before(first) || done.After(day) {
			continue
		}
		out = append(out, t)
	}
	slices.SortStableFunc(out, func(a, b taskfiles.Task) int {
		return cmp.Or(cmp.Compare(b.Done, a.Done), cmp.Compare(a.ID, b.ID))
	})
	return out
}

func archivedCard(t taskfiles.Task) Card {
	return Card{
		ID:           t.ID,
		Title:        t.Title,
		Status:       t.Status,
		Owners:       orEmpty(t.Owners),
		Testers:      orEmpty(t.Testers),
		Labels:       labelStrings(t.Labels),
		Prefix:       idPrefix(t.ID),
		Added:        t.Added,
		Claimed:      t.Claimed,
		Done:         t.Done,
		PRs:          orEmpty(t.PRs),
		PRRef:        t.PRRef,
		OpenPRs:      []LinkedPR{},
		BlockedBy:    t.BlockedBy,
		BlockedByIDs: []string{},
		Blocks:       []string{},
		References:   []string{},
		ReferencedBy: []string{},
		Archived:     true,
		Path:         t.Path,
	}
}

func problems(in []taskfiles.Problem) []Problem {
	out := make([]Problem, 0, len(in))
	for _, p := range in {
		out = append(out, Problem(p))
	}
	return out
}

func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func buildMaps(parsed taskfiles.Result, cards []Card) []EpicMap {
	maps := make([]EpicMap, 0, len(parsed.Maps))
	for _, m := range parsed.Maps {
		label := taskfiles.Label{Group: "epic", Value: m.Epic}.String()
		counts := map[string]int{}
		for _, c := range cards {
			if slices.Contains(c.Labels, label) {
				counts[c.Column]++
			}
		}
		archived := 0
		for _, t := range parsed.Archived {
			if slices.ContainsFunc(t.Labels, func(l taskfiles.Label) bool { return l.String() == label }) {
				archived++
			}
		}
		maps = append(maps, EpicMap{Epic: m.Epic, Body: m.Body, Counts: counts, Archived: archived})
	}
	return maps
}
