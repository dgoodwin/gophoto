package config

import "testing"
import "strings"

func TestLoadSimpleConfig(t *testing.T) {
	cfg := LoadConfigFile("../gophoto-docker.yml")
	if cfg.AssetsPath != "./public" {
		t.Errorf("Unexpected asset path: '%s' != '%s'", cfg.AssetsPath, "./public")
	}

	if cfg.LocalStorage == nil {
		t.Error("Local storage backend was not parsed")
	}

	if cfg.LocalStorage.Path != "./storage/" {
		t.Errorf("Unexpected storage path: '%s' != '%s'",
			cfg.LocalStorage.Path, "./storage/")
	}

	if !strings.Contains(cfg.Database.Open, "user=postgres") ||
		!strings.Contains(cfg.Database.Open, "dbname=gophoto") {
		t.Errorf("Unexpected database open: '%s'", cfg.Database.Open)
	}
}
