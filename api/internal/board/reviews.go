package board

import (
	"context"
	"errors"
	"maps"
	"path"
	"slices"

	"fuda/internal/taskfiles"
)

type PullRequest struct {
	Number       int
	URL          string
	HeadSHA      string
	ChangedPaths []string
}

type ReviewSource interface {
	OpenPRs(ctx context.Context, base string) ([]PullRequest, error)
	FileAt(ctx context.Context, repoPath, ref string) ([]byte, error)
}

type reviews struct {
	open  map[string][]int
	added map[string]taskfiles.Task
}

func (s *Service) openReviews(ctx context.Context, develop taskfiles.Result) (reviews, error) {
	found := reviews{open: map[string][]int{}, added: map[string]taskfiles.Task{}}
	source, ok := s.source.(ReviewSource)
	if !ok {
		return found, nil
	}
	pulls, err := source.OpenPRs(ctx, s.opts.WorkBranch)
	if err != nil {
		return reviews{}, err
	}

	byPath := map[string]taskfiles.Task{}
	for _, t := range develop.Tasks {
		byPath[t.Path] = t
	}
	tasksDir := path.Join(s.opts.BoardDir, "tasks")
	for _, pr := range pulls {
		for _, p := range pr.ChangedPaths {
			if path.Dir(p) != tasksDir || path.Ext(p) != ".md" {
				continue
			}
			content, err := source.FileAt(ctx, p, pr.HeadSHA)
			if errors.Is(err, ErrNotFound) {
				continue
			}
			if err != nil {
				return reviews{}, err
			}
			after := taskfiles.Parse(map[string][]byte{p: content}, taskfiles.Layout{BoardDir: s.opts.BoardDir})
			if len(after.Tasks) != 1 {
				continue
			}
			var before *taskfiles.Task
			if t, ok := byPath[p]; ok {
				before = &t
			}
			if !Delivers(pr.Number, before, after.Tasks[0]) {
				continue
			}
			delivered := after.Tasks[0]
			found.open[delivered.ID] = append(found.open[delivered.ID], pr.Number)
			if before == nil {
				found.added[delivered.ID] = delivered
			}
		}
	}
	return found, nil
}

func sameReviews(a, b reviews) bool {
	return maps.EqualFunc(a.open, b.open, slices.Equal[[]int]) && maps.EqualFunc(a.added, b.added, func(x, y taskfiles.Task) bool {
		return x.Path == y.Path && x.Title == y.Title && x.Status == y.Status
	})
}
