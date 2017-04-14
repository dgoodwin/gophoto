package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type GophotoConfig struct {
	AssetsPath   string `yaml:"assetsPath"`
	Database     Database
	LocalStorage *LocalStorage `yaml:"localStorage"`
	S3Storage    *S3Storage    `yaml:"s3Storage"`
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
