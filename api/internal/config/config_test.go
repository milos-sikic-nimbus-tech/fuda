package config

import "testing"

func TestSourceTakesAList(t *testing.T) {
	t.Setenv("FUDA_SOURCE", "github,azure")
	t.Setenv("FUDA_GITHUB_CLIENT_ID", "g")
	t.Setenv("FUDA_GITHUB_CLIENT_SECRET", "g")
	t.Setenv("FUDA_AZURE_CLIENT_ID", "a")
	t.Setenv("FUDA_AZURE_CLIENT_SECRET", "a")
	t.Setenv("FUDA_COOKIE_SECRET", "c")
	t.Setenv("FUDA_BASE_URL", "http://fuda.test")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Has(SourceGitHub) || !cfg.Has(SourceAzure) || cfg.Has(SourceLocal) {
		t.Errorf("sources: %v", cfg.Sources)
	}
}

func TestEachSourceNeedsItsOwnSettings(t *testing.T) {
	t.Setenv("FUDA_SOURCE", "github,azure")
	t.Setenv("FUDA_GITHUB_CLIENT_ID", "g")
	t.Setenv("FUDA_GITHUB_CLIENT_SECRET", "g")
	t.Setenv("FUDA_COOKIE_SECRET", "c")
	t.Setenv("FUDA_BASE_URL", "http://fuda.test")
	if _, err := Load(); err == nil {
		t.Error("azure without its client settings passed")
	}
}

func TestLocalStandsAlone(t *testing.T) {
	t.Setenv("FUDA_SOURCE", "local,github")
	t.Setenv("FUDA_LOCAL_PATH", "/tmp/x")
	if _, err := Load(); err == nil {
		t.Error("local mixed with github passed")
	}
}
