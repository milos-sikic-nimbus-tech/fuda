package board

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"fuda/internal/taskfiles"
)

const cacheFile = "snapshot.gob"

type cachedFiles struct {
	Heads   map[string]string
	Files   map[string][]byte
	SavedAt time.Time
}

func (s *Service) saveCache(heads map[string]string, files map[string][]byte) error {
	if s.opts.CacheDir == "" {
		return nil
	}
	if err := os.MkdirAll(s.opts.CacheDir, 0o750); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(s.opts.CacheDir, cacheFile+".*")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	err = gob.NewEncoder(tmp).Encode(cachedFiles{Heads: heads, Files: files, SavedAt: s.now()})
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmp.Name(), filepath.Join(s.opts.CacheDir, cacheFile))
}

func (s *Service) Restore() error {
	if s.opts.CacheDir == "" {
		return nil
	}
	f, err := os.Open(filepath.Join(s.opts.CacheDir, cacheFile))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	var cached cachedFiles
	if err := gob.NewDecoder(f).Decode(&cached); err != nil {
		return fmt.Errorf("read %s: %w", f.Name(), err)
	}
	if len(cached.Files) == 0 {
		return nil
	}
	snap := s.build(taskfiles.Parse(cached.Files, s.layout()), cached.Files, cached.Heads, prsByTask{})
	s.snapshot.Store(snap)
	s.markSyncedAt(cached.Heads, cached.SavedAt)
	return nil
}
