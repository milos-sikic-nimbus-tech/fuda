package app

import (
	"context"
	"log/slog"
	"path"
	"path/filepath"
	"strings"

	"fuda/internal/board"
	"fuda/internal/config"
	"fuda/internal/source/azure"
	"fuda/internal/source/github"
	"fuda/internal/source/local"
)

const (
	docsRoot   = "docs"
	boardDir   = "docs/board"
	workBranch = "develop"
	prodBranch = "main"
)

type hostLinks struct {
	code   func(string) string
	pr     string
	origin board.Origin
}

func newSource(cfg config.Config, id board.BoardID) (board.Source, hostLinks, error) {
	switch id.Host {
	case "local":
		if id.Repo != filepath.Base(filepath.Clean(cfg.LocalPath)) {
			return nil, hostLinks{}, board.ErrNotFound
		}
		return local.New(cfg.LocalPath, docsRoot, workBranch), hostLinks{origin: board.Origin{Host: "local", Repo: id.Repo}}, nil
	case "github":
		if cfg.Source != config.SourceGitHub {
			return nil, hostLinks{}, board.ErrNotFound
		}
		gh := github.New(id.Repo, docsRoot)
		return gh, hostLinks{
			code:   gh.CodeURL(workBranch),
			pr:     gh.PRLink(),
			origin: board.Origin{Host: "github", Repo: id.Repo, URL: "https://github.com/" + id.Repo},
		}, nil
	case "azure":
		org, project, name, ok := splitAzure(id.Repo)
		if !ok || (cfg.AzurePAT == "" && cfg.AzureBearer == "") {
			return nil, hostLinks{}, board.ErrNotFound
		}
		repo := azure.Repo{Org: org, Project: project, Name: name}
		az := azure.WithPAT(repo, cfg.AzurePAT, docsRoot)
		if cfg.AzurePAT == "" {
			az = azure.WithBearer(repo, cfg.AzureBearer, docsRoot)
		}
		return az, hostLinks{
			code:   az.CodeURL(workBranch),
			pr:     az.PRLink(),
			origin: board.Origin{Host: "azure", Repo: name, URL: repo.WebURL()},
		}, nil
	}
	return nil, hostLinks{}, board.ErrNotFound
}

func splitAzure(repo string) (org, project, name string, ok bool) {
	parts := strings.Split(repo, "/")
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func envBoard(cfg config.Config) board.BoardID {
	switch cfg.Source {
	case config.SourceAzure:
		return board.BoardID{Host: "azure", Repo: cfg.AzureOrg + "/" + cfg.AzureProject + "/" + cfg.AzureRepo}
	case config.SourceLocal:
		return board.BoardID{Host: "local", Repo: filepath.Base(filepath.Clean(cfg.LocalPath))}
	}
	return board.BoardID{}
}

func listBoards(cfg config.Config) func(context.Context) ([]board.BoardID, error) {
	return func(ctx context.Context) ([]board.BoardID, error) {
		if cfg.Source != config.SourceGitHub {
			return []board.BoardID{envBoard(cfg)}, nil
		}
		repos, err := github.ListBoards(ctx)
		if err != nil {
			return nil, err
		}
		ids := make([]board.BoardID, len(repos))
		for i, repo := range repos {
			ids[i] = board.BoardID{Host: "github", Repo: repo}
		}
		return ids, nil
	}
}

func NewBoards(log *slog.Logger, cfg config.Config) *board.Boards {
	home := envBoard(cfg)
	return board.NewBoards(log, func(id board.BoardID) (*board.Service, error) {
		source, links, err := newSource(cfg, id)
		if err != nil {
			return nil, err
		}
		links.origin.Path = id.Path()
		title := path.Base(id.Repo)
		if id == home && cfg.Title != "" {
			title = cfg.Title
		}
		return board.NewService(source, board.Options{
			Title:      title,
			DocsRoot:   docsRoot,
			BoardDir:   boardDir,
			WorkBranch: workBranch,
			ProdBranch: prodBranch,
			WatchMain:  cfg.WatchMain,
			Cooldown:   cfg.SyncCooldown,
			CacheDir:   filepath.Join(cfg.CacheDir, id.Host, filepath.FromSlash(id.Repo)),
			Logger:     log,
			CodeURL:    links.code,
			PRLink:     links.pr,
			Origin:     links.origin,
		}), nil
	}, listBoards(cfg))
}
