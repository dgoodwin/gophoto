package config

import "testing"

func TestLoadSimpleConfig(t *testing.T) {
	cfg := LoadConfig("../devel-config.yml")
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
}
