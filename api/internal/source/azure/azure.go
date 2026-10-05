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
		base:          "https://dev.azure.com/" + url.PathEscape(repo.Org) + "/" + url.PathEscape(repo.Project) + "/_apis/git/repositories/" + url.PathEscape(repo.Name),
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

func (s *Source) OpenPRs(ctx context.Context, base string) ([]board.PullRequest, error) {
	var pulls struct {
		Value []struct {
			ID                    int `json:"pullRequestId"`
			LastMergeSourceCommit struct {
				CommitID string `json:"commitId"`
			} `json:"lastMergeSourceCommit"`
		} `json:"value"`
	}
	if err := s.getJSON(ctx, "/pullrequests", url.Values{
		"searchCriteria.status":        {"active"},
		"searchCriteria.targetRefName": {"refs/heads/" + base},
		"$top":                         {"100"},
	}, &pulls); err != nil {
		return nil, err
	}
	out := make([]board.PullRequest, 0, len(pulls.Value))
	for _, p := range pulls.Value {
		paths, err := s.changedPaths(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, board.PullRequest{
			Number:       p.ID,
			URL:          strings.ReplaceAll(s.PRLink(), "{n}", fmt.Sprint(p.ID)),
			HeadSHA:      p.LastMergeSourceCommit.CommitID,
			ChangedPaths: paths,
		})
	}
	return out, nil
}

func (s *Source) changedPaths(ctx context.Context, pr int) ([]string, error) {
	var iterations struct {
		Value []struct {
			ID int `json:"id"`
		} `json:"value"`
	}
	if err := s.getJSON(ctx, fmt.Sprintf("/pullrequests/%d/iterations", pr), nil, &iterations); err != nil {
		return nil, err
	}
	if len(iterations.Value) == 0 {
		return nil, nil
	}
	latest := iterations.Value[len(iterations.Value)-1].ID
	var changes struct {
		ChangeEntries []struct {
			ChangeType string `json:"changeType"`
			Item       struct {
				Path string `json:"path"`
			} `json:"item"`
		} `json:"changeEntries"`
	}
	if err := s.getJSON(ctx, fmt.Sprintf("/pullrequests/%d/iterations/%d/changes", pr, latest), url.Values{"$compareTo": {"0"}, "$top": {"2000"}}, &changes); err != nil {
		return nil, err
	}
	var paths []string
	for _, c := range changes.ChangeEntries {
		if !strings.Contains(c.ChangeType, "delete") && c.Item.Path != "" {
			paths = append(paths, strings.TrimPrefix(c.Item.Path, "/"))
		}
	}
	return paths, nil
}

func (s *Source) FileAt(ctx context.Context, repoPath, ref string) ([]byte, error) {
	content, err := s.get(ctx, "/items", url.Values{
		"path":                          {"/" + repoPath},
		"versionDescriptor.version":     {ref},
		"versionDescriptor.versionType": {"commit"},
		"$format":                       {"octetStream"},
	})
	if errors.Is(err, errNotFound) {
		return nil, board.ErrNotFound
	}
	return content, err
}

func (s *Source) getJSON(ctx context.Context, path string, query url.Values, v any) error {
	body, err := s.get(ctx, path, query)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (s *Source) get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("api-version", apiVersion)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.base+path+"?"+query.Encode(), nil)
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
