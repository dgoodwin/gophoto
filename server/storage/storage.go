package storage

import (
	"errors"

	"github.com/dgoodwin/gophoto/config"

	log "github.com/Sirupsen/logrus"
)

func NewStorageBackend(cfg config.GophotoConfig) (StorageBackend, error) {
	log.Infoln("Initializing storage backend: ")
	if cfg.LocalStorage != nil {
		return LocalStorage{Path: cfg.LocalStorage.Path}, nil
	}
	return nil, errors.New("no storage backend defined")
}

type StorageBackend interface {
	ImportFilePath(filepath string) error
}

type LocalStorage struct {
	Path string
}

func (ls LocalStorage) ImportFilePath(filepath string) error {
	return nil
}
