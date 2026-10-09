package azure

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
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

const apiVersion = "7.1"

var errNotFound = errors.New("azure: not found")

type Repo struct {
	Org     string
	Project string
	Name    string
}

func (r Repo) WebURL() string {
	return "https://dev.azure.com/" + url.PathEscape(r.Org) + "/" + url.PathEscape(r.Project) + "/_git/" + url.PathEscape(r.Name)
}

func (r Repo) apiBase() string {
	return "https://dev.azure.com/" + url.PathEscape(r.Org) + "/" + url.PathEscape(r.Project) + "/_apis/git/repositories/" + url.PathEscape(r.Name)
}

type Source struct {
	repo          Repo
	authorization string
	docsRoot      string
	base          string
	client        *http.Client
}

func WithPAT(repo Repo, pat, docsRoot string) *Source {
	return newSource(repo, "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+pat)), docsRoot)
}

func WithBearer(repo Repo, token, docsRoot string) *Source {
	return newSource(repo, "Bearer "+token, docsRoot)
}

func newSource(repo Repo, authorization, docsRoot string) *Source {
	return &Source{
		repo:          repo,
		authorization: authorization,
		docsRoot:      docsRoot,
		base:          repo.apiBase(),
		client:        &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *Source) CodeURL(branch string) func(string) string {
	return func(repoPath string) string {
		return s.repo.WebURL() + "?path=" + url.QueryEscape("/"+repoPath) + "&version=GB" + url.QueryEscape(branch)
	}
}

func (s *Source) PRLink() string {
	return s.repo.WebURL() + "/pullrequest/{n}"
}

func (s *Source) Head(ctx context.Context, branch string) (string, error) {
	var out struct {
		Value []struct {
			Name     string `json:"name"`
			ObjectID string `json:"objectId"`
		} `json:"value"`
	}
	if err := s.getJSON(ctx, "/refs", url.Values{"filter": {"heads/" + branch}}, &out); err != nil {
		return "", err
	}
	for _, ref := range out.Value {
		if ref.Name == "refs/heads/"+branch {
			return ref.ObjectID, nil
		}
	}
	return "", board.ErrBranchMissing
}

func (s *Source) Files(ctx context.Context, branch string) (map[string][]byte, error) {
	body, err := s.get(ctx, "/items", url.Values{
		"path":                          {"/" + s.docsRoot},
		"$format":                       {"zip"},
		"download":                      {"true"},
		"versionDescriptor.version":     {branch},
		"versionDescriptor.versionType": {"branch"},
	})
	if errors.Is(err, errNotFound) {
		return map[string][]byte{}, nil
	}
	if err != nil {
		return nil, err
	}
	archive, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("items zip: %w", err)
	}
	files := map[string][]byte{}
	for _, f := range archive.File {
		name := strings.TrimPrefix(f.Name, "/")
		if f.FileInfo().IsDir() || !strings.HasPrefix(name, s.docsRoot+"/") || !board.WantedFile(name, int64(f.UncompressedSize64)) {
			continue
		}
		content, err := readZipFile(f)
		if err != nil {
			return nil, err
		}
		files[name] = content
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

func (s *Source) OpenPRs(ctx context.Context, repo string) ([]board.PullRequest, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 3 {
		return nil, board.ErrNotFound
	}
	code := Repo{Org: parts[0], Project: parts[1], Name: parts[2]}
	var pulls struct {
		Value []struct {
			ID            int    `json:"pullRequestId"`
			Title         string `json:"title"`
			SourceRefName string `json:"sourceRefName"`
		} `json:"value"`
	}
	err := s.getJSONFrom(ctx, code.apiBase(), "/pullrequests", url.Values{
		"searchCriteria.status": {"active"},
		"$top":                  {"100"},
	}, &pulls)
	if errors.Is(err, errNotFound) {
		return nil, board.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	out := make([]board.PullRequest, len(pulls.Value))
	for i, p := range pulls.Value {
		out[i] = board.PullRequest{
			Number: p.ID,
			URL:    fmt.Sprintf("%s/pullrequest/%d", code.WebURL(), p.ID),
			Title:  p.Title,
			Branch: strings.TrimPrefix(p.SourceRefName, "refs/heads/"),
		}
	}
	return out, nil
}

func (s *Source) getJSON(ctx context.Context, path string, query url.Values, v any) error {
	return s.getJSONFrom(ctx, s.base, path, query, v)
}

func (s *Source) getJSONFrom(ctx context.Context, base, path string, query url.Values, v any) error {
	body, err := s.fetch(ctx, base, path, query)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (s *Source) get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	return s.fetch(ctx, s.base, path, query)
}

func (s *Source) fetch(ctx context.Context, base, path string, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("api-version", apiVersion)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+path+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.authorization)
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
	case res.StatusCode == http.StatusNonAuthoritativeInfo, res.StatusCode == http.StatusUnauthorized:
		return nil, fmt.Errorf("azure %s: authentication failed (%s); check FUDA_AZURE_PAT", path, res.Status)
	case res.StatusCode >= 300:
		return nil, fmt.Errorf("azure %s: %s: %s", path, res.Status, truncate(body))
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
