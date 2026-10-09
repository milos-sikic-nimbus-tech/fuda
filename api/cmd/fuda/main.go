package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fuda/internal/app"
	"fuda/internal/config"
	"fuda/internal/httpapi"
	"fuda/internal/login"
	"fuda/internal/web"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("fuda stopped", "error", err)
		os.Exit(1)
	}
}

func newLogin(log *slog.Logger, cfg config.Config) (httpapi.Login, error) {
	web := login.Config{CookieSecret: cfg.CookieSecret, BaseURL: cfg.BaseURL}
	switch cfg.Source {
	case config.SourceGitHub:
		web.ClientID, web.ClientSecret = cfg.GitHubClientID, cfg.GitHubClientSecret
		return login.NewGitHub(log, web)
	case config.SourceAzure:
		web.ClientID, web.ClientSecret, web.Tenant = cfg.AzureClientID, cfg.AzureClientSecret, cfg.AzureTenant
		return login.NewAzure(log, web)
	}
	return nil, nil
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	hostLogin, err := newLogin(log, cfg)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewHandler(log, app.NewBoards(log, cfg), web.Dist(), hostLogin),
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
