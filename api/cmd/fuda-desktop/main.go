package main

import (
	"cmp"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"fuda/internal/app"
	"fuda/internal/config"
	"fuda/internal/httpapi"
	"fuda/internal/keychain"
	"fuda/internal/login"
	"fuda/internal/web"
)

var (
	version        = "dev"
	githubClientID = ""
	azureClientID  = ""
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	if err := run(log); err != nil {
		log.Error("fuda stopped", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	clientID := githubClientID
	if fromEnv := os.Getenv("FUDA_GITHUB_CLIENT_ID"); fromEnv != "" {
		clientID = fromEnv
	}
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	cfg := config.Config{
		Source:         config.SourceGitHub,
		GitHubClientID: clientID,
		SyncCooldown:   30 * time.Second,
		CacheDir:       filepath.Join(cacheRoot, "fuda"),
	}

	deviceConfig := login.DeviceConfig{
		Host:         "github",
		ClientID:     clientID,
		ClientSecret: os.Getenv("FUDA_GITHUB_CLIENT_SECRET"),
	}
	if os.Getenv("FUDA_SOURCE") == string(config.SourceAzure) {
		cfg.Source = config.SourceAzure
		deviceConfig = login.DeviceConfig{
			Host:     "azure",
			ClientID: cmp.Or(os.Getenv("FUDA_AZURE_CLIENT_ID"), azureClientID),
			Tenant:   cmp.Or(os.Getenv("FUDA_AZURE_TENANT"), "organizations"),
		}
	}
	var desktop *application.App
	deviceConfig.Store = keychain.New(deviceConfig.Host)
	deviceConfig.Open = func(url string) error { return desktop.Browser.OpenURL(url) }
	device := login.NewDevice(log, deviceConfig)
	desktop = application.New(application.Options{
		Name: "fuda",
		Assets: application.AssetOptions{
			Handler: httpapi.NewHandler(log, app.NewBoards(log, cfg), web.Dist(), device),
		},
	})

	menu := desktop.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)
	help := menu.AddSubmenu("Help")
	help.Add("Check for updates…").OnClick(func(*application.Context) {
		go checkForUpdates(desktop, log, true)
	})
	desktop.Menu.Set(menu)

	desktop.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "fuda",
		Width:     1280,
		Height:    820,
		MinWidth:  720,
		MinHeight: 480,
		URL:       "/",
	})
	go checkForUpdates(desktop, log, false)
	return desktop.Run()
}
