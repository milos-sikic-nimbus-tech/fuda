package board

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"fuda/internal/markdown"
	"fuda/internal/taskfiles"
)

const syncTimeout = 2 * time.Minute

func (s *Service) Sync(ctx context.Context) error {
	_, err, _ := s.group.Do("sync", func() (any, error) {
		ctx, cancel := context.WithTimeout(ctx, syncTimeout)
		defer cancel()
		return nil, s.sync(ctx)
	})
	return err
}

func (s *Service) Poll(ctx context.Context) error {
	head, err := s.source.Head(ctx, s.opts.WorkBranch)
	if err != nil {
		return err
	}
	snap := s.snapshot.Load()
	if snap == nil || snap.heads[s.opts.WorkBranch] != head {
		return s.Sync(ctx)
	}
	if s.refreshDue() {
		go func() { _ = s.Sync(context.WithoutCancel(ctx)) }()
	}
	return nil
}

func (s *Service) refreshDue() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.now().Sub(s.status.LastAttempt) >= s.opts.Cooldown
}

func (s *Service) RequestSync(ctx context.Context) bool {
	s.mu.Lock()
	now := s.now()
	if !s.lastRequest.IsZero() && now.Sub(s.lastRequest) < s.opts.Cooldown {
		s.mu.Unlock()
		return false
	}
	s.lastRequest = now
	s.mu.Unlock()

	go func() { _ = s.Sync(context.WithoutCancel(ctx)) }()
	return true
}

func (s *Service) sync(ctx context.Context) error {
	err := s.load(ctx)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.LastAttempt = s.now()
	if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrForbidden) {
		return err
	}
	if err != nil {
		s.status.LastError = err.Error()
		return err
	}
	s.status.LastError = ""
	return nil
}

func (s *Service) load(ctx context.Context) error {
	head, err := s.source.Head(ctx, s.opts.WorkBranch)
	if err != nil {
		return fmt.Errorf("%s: %w", s.opts.WorkBranch, err)
	}
	heads := map[string]string{s.opts.WorkBranch: head}

	current := s.snapshot.Load()
	if current != nil && maps.Equal(current.heads, heads) {
		found, err := s.openReviews(ctx, current.parsed)
		if err != nil {
			s.markSynced(heads)
			return fmt.Errorf("open pull requests: %w", err)
		}
		if !sameReviews(current.reviews, found) {
			current = s.build(current.parsed, current.files, heads, found)
			s.snapshot.Store(current)
		}
		s.markSynced(heads)
		return nil
	}

	files, err := s.source.Files(ctx, s.opts.WorkBranch)
	if err != nil {
		return fmt.Errorf("%s: %w", s.opts.WorkBranch, err)
	}
	parsed := taskfiles.Parse(files, s.layout())
	if err := s.saveCache(heads, files); err != nil {
		s.opts.Logger.Warn("the disk cache was not updated", "error", err)
	}

	found, reviewErr := s.openReviews(ctx, parsed)
	if reviewErr != nil {
		found = prsByTask{}
	}
	next := s.build(parsed, files, heads, found)
	s.snapshot.Store(next)
	s.markSynced(heads)
	if reviewErr != nil {
		return fmt.Errorf("open pull requests: %w", reviewErr)
	}
	return nil
}

func (s *Service) layout() taskfiles.Layout {
	return taskfiles.Layout{BoardDir: s.opts.BoardDir}
}

func (s *Service) build(parsed taskfiles.Result, files map[string][]byte, heads map[string]string, found prsByTask) *snapshot {
	docs := map[string][]byte{}
	docSet := map[string]bool{}
	assets := map[string][]byte{}
	assetSet := map[string]bool{}
	for p, content := range files {
		if !strings.HasPrefix(p, s.opts.DocsRoot+"/") {
			continue
		}
		if strings.HasSuffix(p, ".md") {
			docs[p] = content
			docSet[p] = true
		} else if _, ok := AssetType(p); ok {
			assets[p] = content
			assetSet[p] = true
		}
	}
	byID := map[string]taskfiles.Task{}
	taskPaths := map[string]string{}
	for _, t := range slices.Concat(parsed.Tasks, parsed.Archived) {
		byID[t.ID] = t
		taskPaths[t.Path] = t.ID
	}

	return &snapshot{
		board: Build(Inputs{
			Parsed:  parsed,
			OpenPRs: found,
		}),
		parsed:  parsed,
		files:   files,
		reviews: found,
		docs:    docs,
		assets:  assets,
		byID:    byID,
		heads:   heads,
		links: markdown.Links{
			Prefix:    s.opts.Origin.Path,
			DocsRoot:  s.opts.DocsRoot,
			Docs:      docSet,
			Assets:    assetSet,
			TaskPaths: taskPaths,
			CodeURL:   s.opts.CodeURL,
		},
	}
}

func (s *Service) markSynced(heads map[string]string) {
	s.markSyncedAt(heads, s.now())
}

func (s *Service) markSyncedAt(heads map[string]string, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Branch = BranchStatus{Branch: s.opts.WorkBranch, SHA: heads[s.opts.WorkBranch], SyncedAt: now}
}
