package board

import (
	"context"
	"time"
)

type Access interface {
	CanWrite(ctx context.Context) (bool, error)
}

const accessTTL = time.Minute

type knownAccess struct {
	canWrite bool
	checked  time.Time
}

func (s *Service) CanWrite(ctx context.Context) (bool, error) {
	if _, ok := s.source.(Writer); !ok {
		return false, nil
	}
	access, ok := s.source.(Access)
	if !ok {
		return true, nil
	}
	token := TokenFrom(ctx)
	s.mu.Lock()
	known, found := s.access[token]
	s.mu.Unlock()
	if found && s.now().Sub(known.checked) < accessTTL {
		return known.canWrite, nil
	}
	canWrite, err := access.CanWrite(ctx)
	if err != nil {
		return false, err
	}
	s.mu.Lock()
	s.access[token] = knownAccess{canWrite: canWrite, checked: s.now()}
	s.mu.Unlock()
	return canWrite, nil
}

func (s *Service) requireWriter(ctx context.Context) (Writer, error) {
	canWrite, err := s.CanWrite(ctx)
	if err != nil {
		return nil, err
	}
	if !canWrite {
		return nil, ErrForbidden
	}
	return s.source.(Writer), nil
}
