package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fuda/internal/board"
	"fuda/internal/config"
	"fuda/internal/httpapi"
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
	code func(string) string
	pr   string
}

func newSource(cfg config.Config) (board.Source, hostLinks, error) {
	switch cfg.Source {
	case config.SourceLocal:
		return local.New(cfg.LocalPath, docsRoot, workBranch), hostLinks{}, nil
	case config.SourceGitHub:
		gh := github.New(cfg.GitHubRepo, cfg.GitHubToken, docsRoot)
		return gh, hostLinks{code: gh.CodeURL(workBranch), pr: gh.PRLink()}, nil
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
		Title:        cfg.Title,
		DocsRoot:     docsRoot,
		BoardDir:     boardDir,
		WorkBranch:   workBranch,
		ProdBranch:   prodBranch,
		WatchMain:    cfg.WatchMain,
		ArchiveAfter: time.Duration(cfg.ArchiveAfterDays) * 24 * time.Hour,
		Cooldown:     cfg.SyncCooldown,
		CodeURL:      links.code,
		PRLink:       links.pr,
	})
	go func() {
		if err := service.Sync(ctx); err != nil {
			log.Error("first sync failed; serving without data until a sync succeeds", "error", err)
		}
		service.Run(ctx, cfg.SyncInterval)
	}()

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewHandler(log, service, web.Dist(), httpapi.Credentials{User: cfg.AuthUser, Password: cfg.AuthPassword}, cfg.WebhookSecret),
		ReadHeaderTimeout: 10 * time.Second,
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
