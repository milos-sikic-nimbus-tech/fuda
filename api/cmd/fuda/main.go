package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"fuda/internal/board"
	"fuda/internal/config"
	"fuda/internal/httpapi"
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

func newSource(cfg config.Config) (board.Source, hostLinks, error) {
	switch cfg.Source {
	case config.SourceLocal:
		return local.New(cfg.LocalPath, docsRoot, workBranch), hostLinks{origin: board.Origin{Host: "local", Repo: filepath.Base(filepath.Clean(cfg.LocalPath))}}, nil
	case config.SourceGitHub:
		gh := github.New(cfg.GitHubRepo, cfg.GitHubToken, docsRoot)
		return gh, hostLinks{
			code:   gh.CodeURL(workBranch),
			pr:     gh.PRLink(),
			origin: board.Origin{Host: "github", Repo: cfg.GitHubRepo, URL: "https://github.com/" + cfg.GitHubRepo},
		}, nil
	case config.SourceAzure:
		repo := azure.Repo{Org: cfg.AzureOrg, Project: cfg.AzureProject, Name: cfg.AzureRepo}
		az := azure.WithPAT(repo, cfg.AzurePAT, docsRoot)
		if cfg.AzurePAT == "" {
			az = azure.WithBearer(repo, cfg.AzureBearer, docsRoot)
		}
		return az, hostLinks{
			code:   az.CodeURL(workBranch),
			pr:     az.PRLink(),
			origin: board.Origin{Host: "azure", Repo: cfg.AzureRepo, URL: repo.WebURL()},
		}, nil
	}
	return nil, hostLinks{}, fmt.Errorf("source %q is not available yet", cfg.Source)
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	source, links, err := newSource(cfg)
	if err != nil {
		return err
	}
	service := board.NewService(source, board.Options{
		Title:      cfg.Title,
		DocsRoot:   docsRoot,
		BoardDir:   boardDir,
		WorkBranch: workBranch,
		ProdBranch: prodBranch,
		WatchMain:  cfg.WatchMain,
		Cooldown:   cfg.SyncCooldown,
		CacheDir:   cfg.CacheDir,
		Logger:     log,
		CodeURL:    links.code,
		PRLink:     links.pr,
		Origin:     links.origin,
	})
	if err := service.Restore(); err != nil {
		log.Warn("the disk cache could not be read; starting empty", "error", err)
	}
	go func() {
		if err := service.Sync(ctx); err != nil {
			log.Error("first sync failed; serving the cached copy, if any, until a sync succeeds", "error", err)
		}
		service.Run(ctx, cfg.SyncInterval)
	}()

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewHandler(log, service, web.Dist(), httpapi.Credentials{User: cfg.AuthUser, Password: cfg.AuthPassword}, cfg.WebhookSecret),
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
