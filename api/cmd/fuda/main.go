package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"fuda/internal/board"
	"fuda/internal/config"
	"fuda/internal/httpapi"
	"fuda/internal/login"
	"fuda/internal/source/azure"
	"fuda/internal/source/github"
	"fuda/internal/source/local"
	"fuda/internal/web"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("fuda stopped", "error", err)
		os.Exit(1)
	}
}

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

func newBoards(log *slog.Logger, cfg config.Config) *board.Boards {
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

func newGitHubLogin(log *slog.Logger, cfg config.Config) (httpapi.GitHubLogin, error) {
	if cfg.Source != config.SourceGitHub {
		return nil, nil
	}
	return login.NewGitHub(log, login.Config{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		CookieSecret: cfg.CookieSecret,
		BaseURL:      cfg.BaseURL,
	})
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	githubLogin, err := newGitHubLogin(log, cfg)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewHandler(log, newBoards(log, cfg), web.Dist(), githubLogin),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	errs := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", cfg.Addr, "source", cfg.Source)
		errs <- server.ListenAndServe()
	}()

	select {
	case err := <-errs:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
	}

	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(shutdown)
}
