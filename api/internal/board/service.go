package board

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path"
	"regexp"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"

	"fuda/internal/markdown"
	"fuda/internal/taskfiles"
)

type Source interface {
	Head(ctx context.Context, branch string) (string, error)
	Files(ctx context.Context, branch string) (map[string][]byte, error)
}

var ErrBranchMissing = errors.New("branch not found")

var ErrNotFound = errors.New("not found")

type Options struct {
	Title        string
	DocsRoot     string
	BoardDir     string
	WorkBranch   string
	ProdBranch   string
	WatchMain    bool
	ArchiveAfter time.Duration
	Cooldown     time.Duration
	CodeURL      func(repoPath string) string
	PRLink       string
	Origin       Origin
}

type Origin struct {
	Host string `json:"host"`
	Repo string `json:"repo"`
	URL  string `json:"url,omitempty"`
}

type Service struct {
	source Source
	opts   Options
	now    func() time.Time

	snapshot atomic.Pointer[snapshot]
	group    singleflight.Group

	mu          sync.Mutex
	lastRequest time.Time
	pending     bool
	status      SyncStatus
}

type SyncStatus struct {
	Develop     BranchStatus  `json:"develop"`
	Main        *BranchStatus `json:"main,omitempty"`
	LastAttempt time.Time     `json:"lastAttempt"`
	LastError   string        `json:"lastError,omitempty"`
}

type BranchStatus struct {
	Branch   string    `json:"branch"`
	SHA      string    `json:"sha"`
	SyncedAt time.Time `json:"syncedAt"`
	NotYet   bool      `json:"notYet,omitempty"`
}

type snapshot struct {
	board   Board
	develop taskfiles.Result
	main    *taskfiles.Result
	files   map[string][]byte
	reviews reviews
	docs    map[string][]byte
	assets  map[string][]byte
	links   markdown.Links
	byID    map[string]taskfiles.Task
	heads   map[string]string
}

func NewService(source Source, opts Options) *Service {
	return &Service{source: source, opts: opts, now: time.Now}
}

func (s *Service) Sync(ctx context.Context) error {
	_, err, _ := s.group.Do("sync", func() (any, error) {
		return nil, s.sync(ctx)
	})
	return err
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

func (s *Service) NotifyChange(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pending {
		return
	}
	background := context.WithoutCancel(ctx)
	now := s.now()
	wait := s.opts.Cooldown - now.Sub(s.lastRequest)
	if s.lastRequest.IsZero() || wait <= 0 {
		s.lastRequest = now
		go func() { _ = s.Sync(background) }()
		return
	}
	s.pending = true
	time.AfterFunc(wait, func() {
		s.mu.Lock()
		s.pending = false
		s.lastRequest = s.now()
		s.mu.Unlock()
		_ = s.Sync(background)
	})
}

func (s *Service) sync(ctx context.Context) error {
	err := s.load(ctx)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.LastAttempt = s.now()
	if err != nil {
		s.status.LastError = err.Error()
		return err
	}
	s.status.LastError = ""
	return nil
}

func (s *Service) load(ctx context.Context) error {
	developHead, err := s.source.Head(ctx, s.opts.WorkBranch)
	if err != nil {
		return fmt.Errorf("%s: %w", s.opts.WorkBranch, err)
	}
	heads := map[string]string{s.opts.WorkBranch: developHead}
	if s.opts.WatchMain {
		mainHead, err := s.source.Head(ctx, s.opts.ProdBranch)
		switch {
		case errors.Is(err, ErrBranchMissing):
		case err != nil:
			return fmt.Errorf("%s: %w", s.opts.ProdBranch, err)
		default:
			heads[s.opts.ProdBranch] = mainHead
		}
	}

	current := s.snapshot.Load()
	if current != nil && maps.Equal(current.heads, heads) {
		found, err := s.openReviews(ctx, current.develop)
		if err != nil {
			s.markSynced(heads, current)
			return fmt.Errorf("open pull requests: %w", err)
		}
		if !sameReviews(current.reviews, found) {
			current = s.build(current.develop, current.main, current.files, heads, found)
			s.snapshot.Store(current)
		}
		s.markSynced(heads, current)
		return nil
	}

	developFiles, err := s.source.Files(ctx, s.opts.WorkBranch)
	if err != nil {
		return fmt.Errorf("%s: %w", s.opts.WorkBranch, err)
	}
	layout := taskfiles.Layout{BoardDir: s.opts.BoardDir}
	develop := taskfiles.Parse(developFiles, layout)

	var main *taskfiles.Result
	if _, ok := heads[s.opts.ProdBranch]; ok && s.opts.WatchMain {
		mainFiles, err := s.source.Files(ctx, s.opts.ProdBranch)
		if err != nil {
			return fmt.Errorf("%s: %w", s.opts.ProdBranch, err)
		}
		parsed := taskfiles.Parse(mainFiles, layout)
		if len(parsed.Tasks) > 0 || len(parsed.Archived) > 0 {
			main = &parsed
		}
	}

	found, reviewErr := s.openReviews(ctx, develop)
	if reviewErr != nil {
		found = reviews{}
	}
	next := s.build(develop, main, developFiles, heads, found)
	s.snapshot.Store(next)
	s.markSynced(heads, next)
	if reviewErr != nil {
		return fmt.Errorf("open pull requests: %w", reviewErr)
	}
	return nil
}

func (s *Service) build(develop taskfiles.Result, main *taskfiles.Result, files map[string][]byte, heads map[string]string, found reviews) *snapshot {
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
	for _, t := range slices.Concat(develop.Tasks, develop.Archived) {
		byID[t.ID] = t
		taskPaths[t.Path] = t.ID
	}
	for id, t := range found.added {
		if _, exists := byID[id]; !exists {
			byID[id] = t
		}
	}

	return &snapshot{
		board: Build(Inputs{
			Develop:      develop,
			Main:         main,
			OpenPRs:      found.open,
			AddedInPRs:   slices.Collect(maps.Values(found.added)),
			ArchiveAfter: s.opts.ArchiveAfter,
			Now:          s.now(),
		}),
		develop: develop,
		main:    main,
		files:   files,
		reviews: found,
		docs:    docs,
		assets:  assets,
		byID:    byID,
		heads:   heads,
		links: markdown.Links{
			DocsRoot:  s.opts.DocsRoot,
			Docs:      docSet,
			Assets:    assetSet,
			TaskPaths: taskPaths,
			CodeURL:   s.opts.CodeURL,
		},
	}
}

func (s *Service) markSynced(heads map[string]string, snap *snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	s.status.Develop = BranchStatus{Branch: s.opts.WorkBranch, SHA: heads[s.opts.WorkBranch], SyncedAt: now}
	if !s.opts.WatchMain {
		s.status.Main = nil
		return
	}
	s.status.Main = &BranchStatus{Branch: s.opts.ProdBranch, SHA: heads[s.opts.ProdBranch], SyncedAt: now, NotYet: snap.main == nil}
}

type BoardView struct {
	Board
	Title  string     `json:"title"`
	PRLink string     `json:"prLink,omitempty"`
	Origin Origin     `json:"origin"`
	Sync   SyncStatus `json:"sync"`
}

func (s *Service) Board() (BoardView, bool) {
	snap := s.snapshot.Load()
	if snap == nil {
		return BoardView{Title: s.opts.Title, PRLink: s.opts.PRLink, Origin: s.opts.Origin, Sync: s.Status()}, false
	}
	return BoardView{Title: s.opts.Title, PRLink: s.opts.PRLink, Origin: s.opts.Origin, Board: snap.board, Sync: s.Status()}, true
}

func (s *Service) Status() SyncStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	status := s.status
	if status.Main != nil {
		main := *status.Main
		status.Main = &main
	}
	return status
}

type TaskView struct {
	Card
	Custom []taskfiles.Field `json:"custom"`
	HTML   string            `json:"html"`
}

func (s *Service) Task(id string) (TaskView, error) {
	snap := s.snapshot.Load()
	if snap == nil {
		return TaskView{}, ErrNotFound
	}
	t, ok := snap.byID[id]
	if !ok {
		return TaskView{}, ErrNotFound
	}
	html, err := markdown.Render([]byte(t.Body), t.Path, snap.links)
	if err != nil {
		return TaskView{}, err
	}
	view := TaskView{Custom: orEmpty(t.Custom), HTML: html}
	if card, found := snap.cardByID(id); found {
		view.Card = card
	} else {
		view.Card = archivedCard(t)
	}
	return view, nil
}

func (snap *snapshot) cardByID(id string) (Card, bool) {
	for _, c := range snap.board.Cards {
		if c.ID == id {
			return c, true
		}
	}
	return Card{}, false
}

func (s *Service) Search(query string) []string {
	snap := s.snapshot.Load()
	if snap == nil {
		return []string{}
	}
	return Search(snap.develop.Tasks, query)
}

func (s *Service) Archive() []Card {
	snap := s.snapshot.Load()
	if snap == nil {
		return []Card{}
	}
	out := make([]Card, 0, len(snap.develop.Archived))
	for _, t := range snap.develop.Archived {
		out = append(out, archivedCard(t))
	}
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
		OpenPRs:      []int{},
		BlockedBy:    t.BlockedBy,
		BlockedByIDs: []string{},
		Blocks:       []string{},
		References:   []string{},
		ReferencedBy: []string{},
		Path:         t.Path,
	}
}

var markdownLink = regexp.MustCompile(`\]\(([^)\s#]+)`)

func linksTo(t taskfiles.Task, repoPath string) bool {
	for _, m := range markdownLink.FindAllStringSubmatch(t.Body, -1) {
		if path.Clean(path.Join(path.Dir(t.Path), m[1])) == repoPath {
			return true
		}
	}
	return false
}

type DocView struct {
	Path      string   `json:"path"`
	Title     string   `json:"title"`
	HTML      string   `json:"html"`
	Backlinks []string `json:"backlinks"`
}

func (s *Service) Doc(relative string) (DocView, error) {
	snap := s.snapshot.Load()
	if snap == nil {
		return DocView{}, ErrNotFound
	}
	repoPath := path.Clean(path.Join(s.opts.DocsRoot, relative))
	content, ok := snap.docs[repoPath]
	if !ok || !strings.HasPrefix(repoPath, s.opts.DocsRoot+"/") {
		return DocView{}, ErrNotFound
	}
	html, err := markdown.Render(content, repoPath, snap.links)
	if err != nil {
		return DocView{}, err
	}
	name := path.Base(repoPath)
	backlinks := []string{}
	for _, t := range snap.develop.Tasks {
		if linksTo(t, repoPath) {
			backlinks = append(backlinks, t.ID)
		}
	}
	return DocView{
		Path:      strings.TrimPrefix(repoPath, s.opts.DocsRoot+"/"),
		Title:     markdown.Title(content, name),
		HTML:      html,
		Backlinks: backlinks,
	}, nil
}

func (s *Service) Asset(repoPath string) ([]byte, string, error) {
	snap := s.snapshot.Load()
	if snap == nil {
		return nil, "", ErrNotFound
	}
	content, ok := snap.assets[path.Clean(repoPath)]
	if !ok {
		return nil, "", ErrNotFound
	}
	contentType, _ := AssetType(repoPath)
	return content, contentType, nil
}
