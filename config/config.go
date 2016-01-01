package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type GophotoConfig struct {
	AssetsPath string `yaml:"assetspath"`
	Storage    Storage
}

type Storage struct {
	Backend string
	Path    string
}

func LoadConfig(path string) GophotoConfig {
	yamlFile, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	var cfg GophotoConfig

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
