package board

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type BoardID struct {
	Host string
	Repo string
}

var segment = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._ -]*$`)

func (id BoardID) Path() string {
	return (&url.URL{Path: "/" + id.Host + "/" + id.Repo}).EscapedPath()
}

func (id BoardID) valid() bool {
	if id.Host == "" || id.Repo == "" {
		return false
	}
	for part := range strings.SplitSeq(id.Repo, "/") {
		if !segment.MatchString(part) || part == ".." {
			return false
		}
	}
	return true
}

type Boards struct {
	ctx      context.Context
	interval time.Duration
	open     func(BoardID) (*Service, error)
	log      *slog.Logger

	opening  singleflight.Group
	mu       sync.Mutex
	services map[BoardID]*Service
}

func NewBoards(ctx context.Context, log *slog.Logger, interval time.Duration, open func(BoardID) (*Service, error)) *Boards {
	return &Boards{ctx: ctx, interval: interval, open: open, log: log, services: map[BoardID]*Service{}}
}

func (b *Boards) Get(id BoardID) (*Service, error) {
	if !id.valid() {
		return nil, ErrNotFound
	}
	if service, ok := b.cached(id); ok {
		return service, nil
	}
	value, err, _ := b.opening.Do(id.Path(), func() (any, error) {
		if service, ok := b.cached(id); ok {
			return service, nil
		}
		return b.start(id)
	})
	if err != nil {
		return nil, err
	}
	return value.(*Service), nil
}

func (b *Boards) NotifyHost(ctx context.Context, host string) {
	b.mu.Lock()
	var open []*Service
	for id, service := range b.services {
		if id.Host == host {
			open = append(open, service)
		}
	}
	b.mu.Unlock()
	for _, service := range open {
		service.NotifyChange(ctx)
	}
}

func (b *Boards) cached(id BoardID) (*Service, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	service, ok := b.services[id]
	return service, ok
}

func (b *Boards) start(id BoardID) (*Service, error) {
	service, err := b.open(id)
	if err != nil {
		return nil, err
	}
	if err := service.Restore(); err != nil {
		b.log.Warn("the disk cache could not be read; starting empty", "board", id.Path(), "error", err)
	}
	_, restored := service.Board()
	if !restored {
		if err := service.Sync(b.ctx); errors.Is(err, ErrBranchMissing) {
			return nil, ErrNotFound
		}
	}
	go func() {
		if restored {
			_ = service.Sync(b.ctx)
		}
		service.Run(b.ctx, b.interval)
	}()
	b.mu.Lock()
	b.services[id] = service
	b.mu.Unlock()
	return service, nil
}
