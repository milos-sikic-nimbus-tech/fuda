package board

import "context"

type Viewer interface {
	ViewerLogin(ctx context.Context) (string, error)
}

type Me struct {
	Login    string `json:"login"`
	PersonID string `json:"personId,omitempty"`
}

func (s *Service) Me(ctx context.Context) (*Me, error) {
	viewer, ok := s.source.(Viewer)
	if !ok {
		return nil, nil
	}
	token := tokenDigest(ctx)
	s.mu.Lock()
	login, found := s.logins[token]
	s.mu.Unlock()
	if !found {
		var err error
		if login, err = viewer.ViewerLogin(ctx); err != nil {
			return nil, err
		}
		s.mu.Lock()
		s.logins[token] = login
		s.mu.Unlock()
	}
	if login == "" {
		return nil, nil
	}
	me := &Me{Login: login}
	if snap := s.snapshot.Load(); snap != nil {
		me.PersonID = newPeople(snap.parsed.Config.People, snap.parsed.Tasks).ofGitHub(login)
	}
	return me, nil
}
