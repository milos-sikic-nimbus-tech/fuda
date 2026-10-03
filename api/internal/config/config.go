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

	AuthUser     string `env:"FUDA_AUTH_USER" envDefault:"fuda"`
	AuthPassword string `env:"FUDA_AUTH_PASSWORD"`

	WebhookSecret string `env:"FUDA_WEBHOOK_SECRET"`

	Source    Source `env:"FUDA_SOURCE" envDefault:"local"`
	LocalPath string `env:"FUDA_LOCAL_PATH"`

	GitHubRepo  string `env:"FUDA_GITHUB_REPO"`
	GitHubToken string `env:"FUDA_GITHUB_TOKEN"`

	AzureOrg     string `env:"FUDA_AZURE_ORG"`
	AzureProject string `env:"FUDA_AZURE_PROJECT"`
	AzureRepo    string `env:"FUDA_AZURE_REPO"`
	AzurePAT     string `env:"FUDA_AZURE_PAT"`
	AzureBearer  string `env:"FUDA_AZURE_BEARER"`

	WatchMain        bool          `env:"FUDA_WATCH_MAIN" envDefault:"false"`
	ArchiveAfterDays int           `env:"FUDA_ARCHIVE_AFTER_DAYS" envDefault:"15"`
	SyncInterval     time.Duration `env:"FUDA_SYNC_INTERVAL" envDefault:"3m"`
	SyncCooldown     time.Duration `env:"FUDA_SYNC_COOLDOWN" envDefault:"30s"`
	CacheDir         string        `env:"FUDA_CACHE_DIR" envDefault:"./.fuda-cache"`
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
		if c.GitHubRepo == "" || c.GitHubToken == "" {
			return fmt.Errorf("FUDA_SOURCE=github needs FUDA_GITHUB_REPO and FUDA_GITHUB_TOKEN")
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
