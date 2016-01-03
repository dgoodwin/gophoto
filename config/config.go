package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type GophotoConfig struct {
	AssetsPath string `yaml:"assetspath"`
	Storage    Storage
	Database   Database
}

type Storage struct {
	Backend string
	Path    string
}

type Database struct {
	Open string
}

func LoadConfig(configPath string) GophotoConfig {
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
