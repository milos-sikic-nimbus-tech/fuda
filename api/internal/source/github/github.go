package github

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"fuda/internal/board"
)

const apiBase = "https://api.github.com"

var errNotFound = errors.New("github: not found")

type Source struct {
	repo     string
	token    string
	docsRoot string
	base     string
	client   *http.Client
}

func New(repo, token, docsRoot string) *Source {
	return &Source{
		repo:     repo,
		token:    token,
		docsRoot: docsRoot,
		base:     apiBase,
		client:   &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *Source) CodeURL(branch string) func(string) string {
	return func(repoPath string) string {
		return "https://github.com/" + s.repo + "/blob/" + branch + "/" + repoPath
	}
}

func (s *Source) PRLink() string {
	return "https://github.com/" + s.repo + "/pull/{n}"
}

func (s *Source) Head(ctx context.Context, branch string) (string, error) {
	var out struct {
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}
	err := s.getJSON(ctx, "/repos/"+s.repo+"/branches/"+url.PathEscape(branch), &out)
	if errors.Is(err, errNotFound) {
		return "", board.ErrBranchMissing
	}
	if err != nil {
		return "", err
	}
	return out.Commit.SHA, nil
}

func (s *Source) Files(ctx context.Context, branch string) (map[string][]byte, error) {
	body, err := s.get(ctx, "/repos/"+s.repo+"/zipball/"+url.PathEscape(branch), "")
	if err != nil {
		return nil, err
	}
	archive, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("zipball: %w", err)
	}
	files := map[string][]byte{}
	for _, f := range archive.File {
		_, rel, found := strings.Cut(f.Name, "/")
		if !found || f.FileInfo().IsDir() || !strings.HasSuffix(rel, ".md") || !strings.HasPrefix(rel, s.docsRoot+"/") {
			continue
		}
		content, err := readZipFile(f)
		if err != nil {
			return nil, err
		}
		files[rel] = content
	}
	return files, nil
}

func readZipFile(f *zip.File) ([]byte, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return io.ReadAll(r)
}

func (s *Source) OpenPRs(ctx context.Context, base string) ([]board.PullRequest, error) {
	var pulls []struct {
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
		Head    struct {
			SHA string `json:"sha"`
		} `json:"head"`
	}
	if err := s.getJSON(ctx, "/repos/"+s.repo+"/pulls?state=open&per_page=100&base="+url.QueryEscape(base), &pulls); err != nil {
		return nil, err
	}
	out := make([]board.PullRequest, 0, len(pulls))
	for _, p := range pulls {
		var files []struct {
			Filename string `json:"filename"`
			Status   string `json:"status"`
		}
		if err := s.getJSON(ctx, fmt.Sprintf("/repos/%s/pulls/%d/files?per_page=100", s.repo, p.Number), &files); err != nil {
			return nil, err
		}
		pr := board.PullRequest{Number: p.Number, URL: p.HTMLURL, HeadSHA: p.Head.SHA}
		for _, f := range files {
			if f.Status != "removed" {
				pr.ChangedPaths = append(pr.ChangedPaths, f.Filename)
			}
		}
		out = append(out, pr)
	}
	return out, nil
}

func (s *Source) FileAt(ctx context.Context, repoPath, ref string) ([]byte, error) {
	content, err := s.get(ctx, "/repos/"+s.repo+"/contents/"+escapePath(repoPath)+"?ref="+url.QueryEscape(ref), "application/vnd.github.raw+json")
	if errors.Is(err, errNotFound) {
		return nil, board.ErrNotFound
	}
	return content, err
}

func escapePath(p string) string {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func (s *Source) getJSON(ctx context.Context, path string, v any) error {
	body, err := s.get(ctx, path, "application/vnd.github+json")
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (s *Source) get(ctx context.Context, path, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.base+path, nil)
	if err != nil {
		return nil, err
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	switch {
	case res.StatusCode == http.StatusNotFound:
		return nil, errNotFound
	case res.StatusCode >= 300:
		return nil, fmt.Errorf("github %s: %s: %s", path, res.Status, truncate(body))
	}
	return body, nil
}

func truncate(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 200 {
		return s[:200] + "…"
	}
	return s
}
