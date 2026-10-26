package board

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"fuda/internal/taskfiles"
)

type AssignRequest struct {
	TaskID     string   `json:"-"`
	FromOwners []string `json:"seen"`
	Owners     []string `json:"owners"`
}

var ErrUnknownPerson = errors.New("not an id in people.md")

func (s *Service) Assign(ctx context.Context, req AssignRequest) error {
	editor, err := s.requireEditor(ctx)
	if err != nil {
		return err
	}
	snap := s.snapshot.Load()
	if snap == nil {
		return ErrNotFound
	}
	task, ok := snap.byID[req.TaskID]
	if !ok {
		return ErrNotFound
	}
	if task.Archived {
		return ErrNotMovable
	}

	known := newPeople(snap.parsed.Config.People, snap.parsed.Tasks)
	if unknown := known.unknown(req.Owners); len(unknown) > 0 {
		return fmt.Errorf("%w: %s", ErrUnknownPerson, strings.Join(unknown, ", "))
	}
	return s.commit(ctx, editor, task.Path, func(content []byte, current taskfiles.Task) ([]byte, string, error) {
		if !sameNames(known.canonical(current.Owners), known.canonical(req.FromOwners)) {
			who := s.lastEditor(ctx, editor, task.Path)
			now := "no one"
			if owners := known.canonical(current.Owners); len(owners) > 0 {
				now = strings.Join(owners, ", ")
			}
			return nil, "", &ConflictError{Message: fmt.Sprintf("%s set the owners to %s just now", who, now)}
		}
		edited, err := taskfiles.ApplyAssign(content, known.canonical(req.Owners))
		return edited, assignMessage(task.ID, req.Owners), err
	})
}

func assignMessage(id string, owners []string) string {
	if len(owners) == 0 {
		return fmt.Sprintf("Remove the owners of %s", id)
	}
	return fmt.Sprintf("Assign %s to %s", id, strings.Join(owners, ", "))
}

func sameNames(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, name := range a {
		if !slices.Contains(b, name) {
			return false
		}
	}
	return true
}
