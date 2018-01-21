package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"golang.org/x/sys/unix"
)

type GophotoConfig struct {
	AssetsPath   string `yaml:"assetsPath"`
	Database     Database
	LocalStorage *LocalStorage `yaml:"localStorage"`
	S3Storage    *S3Storage    `yaml:"s3Storage"`
	WorkDirPath  string
	APIPort      int `yaml:"apiPort"`
}

type LocalStorage struct {
	Path string
}

// TODO:
type S3Storage struct {
}

type Database struct {
	Open string
}

func LoadConfigFile(configPath string) GophotoConfig {
	yamlFile, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	var cfg GophotoConfig

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	// Expand any environment variables in the database string:
	cfg.Database.Open = os.ExpandEnv(cfg.Database.Open)

	return cfg
}

func validateConfig(cfg *GophotoConfig) error {
	if cfg.WorkDirPath == "" {
		return errors.New("no workDirPath specified in config")
	}
	if _, err := os.Stat(cfg.WorkDirPath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("workDirPath does not exist: %s", cfg.WorkDirPath))
	}
	if unix.Access(cfg.WorkDirPath, unix.W_OK) != nil {
		return errors.New(fmt.Sprintf("workDirPath is not writable: %s", cfg.WorkDirPath))
	}
	if cfg.APIPort == 0 {
		return errors.New(fmt.Sprintf("apiPort not specified in config"))
	}
	return nil
}
