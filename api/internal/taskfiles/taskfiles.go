package taskfiles

import (
	"errors"
	"path"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

type Task struct {
	Path      string
	ID        string
	Title     string
	Status    string
	BlockedBy string
	Owners    []string
	Testers   []string
	Labels    []Label
	Added     string
	Claimed   string
	Done      string
	PRs       []int
	PRRef     string
	Custom    []Field
	Body      string
	Archived  bool
}

type Label struct {
	Group string
	Value string
}

func (l Label) String() string {
	if l.Group == "" {
		return l.Value
	}
	return l.Group + ":" + l.Value
}

type Field struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Stage struct {
	Name       string
	Statuses   []string
	WhenPROpen bool
}

type LabelGroup struct {
	Name   string
	Values []string
}

type Person struct {
	Name    string
	Aliases []string
}

type Config struct {
	Stages      []Stage
	LabelGroups []LabelGroup
	LabelColors map[string]string
	People      []Person
	CodeRepos   []string
}

type Problem struct {
	Path   string
	Reason string
}

type Map struct {
	Epic string
	Path string
	Body string
}

type Result struct {
	Tasks    []Task
	Maps     []Map
	Archived []Task
	Config   Config
	Problems []Problem
}

type Layout struct {
	BoardDir string
}

func (l Layout) tasksDir() string   { return path.Join(l.BoardDir, "tasks") }
func (l Layout) archiveDir() string { return path.Join(l.BoardDir, "archive") }

func Parse(files map[string][]byte, layout Layout) Result {
	var r Result
	r.Config = parseConfig(files, layout, &r.Problems)

	seen := map[string]string{}
	for _, p := range sortedPaths(files) {
		loc, ok := taskLocation(p, layout)
		if ok && path.Base(p) == mapFileName {
			if m, reason := parseEpicMap(p, files[p], loc); reason != "" {
				r.Problems = append(r.Problems, Problem{Path: p, Reason: reason})
			} else if m != nil {
				r.Maps = append(r.Maps, *m)
			}
			continue
		}
		if !ok {
			if nestedUnder(p, layout.tasksDir()) || nestedUnder(p, layout.archiveDir()) {
				r.Problems = append(r.Problems, Problem{Path: p, Reason: "task files must sit directly in tasks/ or archive/, or in one epic folder below them"})
			}
			continue
		}
		task, reason := parseTask(p, files[p], loc.archived)
		if reason == "" {
			if first, dup := seen[task.ID]; dup {
				reason = "duplicate id " + task.ID + " (also in " + first + ")"
			}
		}
		if reason != "" {
			r.Problems = append(r.Problems, Problem{Path: p, Reason: reason})
			continue
		}
		seen[task.ID] = p
		if hasEpicLabel(task.Labels) {
			reason := deprecatedEpicLabels + "; move the task into `tasks/<epic>/`"
			if loc.epic != "" {
				reason = deprecatedEpicLabels + "; remove the label, the folder is the epic"
			}
			r.Problems = append(r.Problems, Problem{Path: p, Reason: reason})
		}
		if loc.epic != "" {
			task.Labels = withEpic(task.Labels, loc.epic)
		}
		if loc.archived {
			r.Archived = append(r.Archived, task)
		} else {
			r.Tasks = append(r.Tasks, task)
		}
	}
	return r
}

const deprecatedEpicLabels = "`epic:` and `theme:` labels are deprecated and will be removed"

func hasEpicLabel(labels []Label) bool {
	return slices.ContainsFunc(labels, isEpicLabel)
}

func isEpicLabel(l Label) bool {
	return l.Group == "epic" || l.Group == "theme"
}

func withEpic(labels []Label, epic string) []Label {
	out := slices.DeleteFunc(slices.Clone(labels), isEpicLabel)
	return append(out, Label{Group: "epic", Value: epic})
}

const (
	mapFileName = "map.md"
	mapLabel    = "wayfinder:map"
)

func parseEpicMap(p string, content []byte, loc location) (*Map, string) {
	if loc.epic == "" || loc.archived {
		return nil, ""
	}
	front, body, err := splitFrontmatter(content)
	if errors.Is(err, errNoFrontmatter) {
		return nil, ""
	}
	if err != nil {
		return nil, err.Error()
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(front, &doc); err != nil {
		return nil, "invalid YAML: " + err.Error()
	}
	fields, err := mappingPairs(&doc)
	if err != nil {
		return nil, err.Error()
	}
	for _, f := range fields {
		if f.key == "label" && f.value.Kind == yaml.ScalarNode && f.value.Value == mapLabel {
			return &Map{Epic: loc.epic, Path: p, Body: string(body)}, ""
		}
	}
	return nil, ""
}

type location struct {
	archived bool
	epic     string
}

func taskLocation(p string, layout Layout) (location, bool) {
	if !strings.HasSuffix(p, ".md") {
		return location{}, false
	}
	dir := path.Dir(p)
	for _, root := range []struct {
		dir      string
		archived bool
	}{{layout.tasksDir(), false}, {layout.archiveDir(), true}} {
		switch {
		case dir == root.dir:
			return location{archived: root.archived}, true
		case path.Dir(dir) == root.dir:
			return location{archived: root.archived, epic: path.Base(dir)}, true
		}
	}
	return location{}, false
}

func nestedUnder(p, dir string) bool {
	return strings.HasSuffix(p, ".md") && strings.HasPrefix(p, dir+"/") && path.Dir(p) != dir
}

func sortedPaths(files map[string][]byte) []string {
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	slices.Sort(paths)
	return paths
}
