package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

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

func TestWorkDir(t *testing.T) {
	cfg, tempDir, err := InitTestConfig()

	if err != nil {
		t.Errorf("Error in test setup: %s", err)
		return
	}

	// Save the valid work dir, some tests change this.
	validWorkDir := cfg.WorkDirPath

	defer os.RemoveAll(tempDir)

	t.Run("workDirPath defined correctly", func(t *testing.T) {
		// the test config should be valid as it is
		if err := validateConfig(cfg); err != nil {
			t.Error("error validating config: ", err)
		}
	})
	t.Run("workDirPath not defined", func(t *testing.T) {
		cfg.WorkDirPath = "" // Blank out the valid work dir
		if err := validateConfig(cfg); err == nil {
			t.Error("no error thrown when workDirPath not defined")
		}
		cfg.WorkDirPath = validWorkDir
	})
	t.Run("workDirPath does not exist", func(t *testing.T) {
		cfg.WorkDirPath = "/tmp/thisdoesntexist"
		if err := validateConfig(cfg); err == nil {
			t.Error("no error thrown when workDirPath does not exist")
		}
		cfg.WorkDirPath = validWorkDir
	})
	t.Run("workDirPath is not writable", func(t *testing.T) {
		// Being lazy, will fail if you're running as root but you're not doing that are you?
		cfg.WorkDirPath = "/bin"
		if err := validateConfig(cfg); err == nil {
			t.Error("no error thrown when workDirPath is not writable")
		}
		cfg.WorkDirPath = validWorkDir
	})
}

// InitTestConfig can be used throughout tests to get a functional config.
// A temporary directory is created to parent all required directories.
// Returns the config, the path to the temporary directory (which caller
// must remember to clean up), and any errors encountered.
func InitTestConfig() (*GophotoConfig, string, error) {
	dir, err := ioutil.TempDir("", "teststorage")
	if err != nil {
		return nil, "", err
	}
	log.Infof("Created temporary test dir: %s", dir)

	// TODO: Create a sub-dir for working directory

	cfg := &GophotoConfig{
		AssetsPath:   "./public", // TODO
		Database:     Database{Open: "user=postgres dbname=gophoto"},
		LocalStorage: &LocalStorage{Path: "./storage"}, // TODO
		WorkDirPath:  dir,
	}

	return cfg, dir, nil
}
