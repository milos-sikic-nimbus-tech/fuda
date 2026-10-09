package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Source string

const (
	SourceLocal  Source = "local"
	SourceGitHub Source = "github"
	SourceAzure  Source = "azure"
)

type Config struct {
	Addr  string `env:"FUDA_ADDR" envDefault:":8080"`
	Title string `env:"FUDA_TITLE"`

	BaseURL      string `env:"FUDA_BASE_URL"`
	CookieSecret string `env:"FUDA_COOKIE_SECRET"`

	Source    Source `env:"FUDA_SOURCE" envDefault:"github"`
	LocalPath string `env:"FUDA_LOCAL_PATH"`

	GitHubClientID     string `env:"FUDA_GITHUB_CLIENT_ID"`
	GitHubClientSecret string `env:"FUDA_GITHUB_CLIENT_SECRET"`

	AzureOrg     string `env:"FUDA_AZURE_ORG"`
	AzureProject string `env:"FUDA_AZURE_PROJECT"`
	AzureRepo    string `env:"FUDA_AZURE_REPO"`
	AzurePAT     string `env:"FUDA_AZURE_PAT"`
	AzureBearer  string `env:"FUDA_AZURE_BEARER"`

	WatchMain    bool          `env:"FUDA_WATCH_MAIN" envDefault:"false"`
	SyncCooldown time.Duration `env:"FUDA_SYNC_COOLDOWN" envDefault:"30s"`
	CacheDir     string        `env:"FUDA_CACHE_DIR" envDefault:"./.fuda-cache"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	return cfg, cfg.validate()
}

func (c Config) validate() error {
	switch c.Source {
	case SourceLocal:
		if c.LocalPath == "" {
			return fmt.Errorf("FUDA_SOURCE=local needs FUDA_LOCAL_PATH")
		}
	case SourceGitHub:
		if c.GitHubClientID == "" || c.GitHubClientSecret == "" || c.CookieSecret == "" || c.BaseURL == "" {
			return fmt.Errorf("FUDA_SOURCE=github needs FUDA_GITHUB_CLIENT_ID, FUDA_GITHUB_CLIENT_SECRET, FUDA_COOKIE_SECRET and FUDA_BASE_URL")
		}
	case SourceAzure:
		if c.AzureOrg == "" || c.AzureProject == "" || c.AzureRepo == "" {
			return fmt.Errorf("FUDA_SOURCE=azure needs FUDA_AZURE_ORG, FUDA_AZURE_PROJECT and FUDA_AZURE_REPO")
		}
		if c.AzurePAT == "" && c.AzureBearer == "" {
			return fmt.Errorf("FUDA_SOURCE=azure needs FUDA_AZURE_PAT (or FUDA_AZURE_BEARER for local runs)")
		}
	default:
		return fmt.Errorf("unknown FUDA_SOURCE %q: use local, github or azure", c.Source)
	}
	return nil
}
