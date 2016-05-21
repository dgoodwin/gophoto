package config

import "testing"
import "strings"

func TestLoadSimpleConfig(t *testing.T) {
	cfg := LoadConfig("../gophoto-docker.yml")
	if cfg.AssetsPath != "./public" {
		t.Errorf("Unexpected asset path: '%s' != '%s'", cfg.AssetsPath, "./public")
	}

	if cfg.Storage.Backend != "fileSystem" {
		t.Errorf("Unexpected storage backend: '%s' != '%s'",
			cfg.Storage.Backend, "fileSystem")
	}

	if cfg.Storage.Path != "./storage/" {
		t.Errorf("Unexpected storage path: '%s' != '%s'",
			cfg.Storage.Path, "./storage/")
	}

	if cfg.ImportPath != "/import/" {
		t.Errorf("Unexpected import path: '%s' != '%s'",
			cfg.ImportPath, "/import/")
	}

	if !strings.Contains(cfg.Database.Open, "user=postgres") ||
		!strings.Contains(cfg.Database.Open, "dbname=gophoto") {
		t.Errorf("Unexpected database open: '%s'", cfg.Database.Open)
	}
}
