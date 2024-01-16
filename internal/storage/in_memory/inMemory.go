package in_memory

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"url_shortener/internal/customError"
	"url_shortener/internal/storage"
)

type inMemoryStorage struct {
	log  *log.Logger
	urls map[string]string
	hash map[string]string
	rwMu *sync.RWMutex
}

func NewMemoryStorage(log *log.Logger) storage.Storage {
	return &inMemoryStorage{
		log:  log,
		urls: make(map[string]string),
		hash: make(map[string]string),
		rwMu: &sync.RWMutex{},
	}
}

func (ms *inMemoryStorage) SaveUrlPair(_ context.Context, shortUrl, originalUrl string) error {
	ms.rwMu.Lock()
	defer ms.rwMu.Unlock()
	if _, ok := ms.urls[originalUrl]; ok {
		return customError.UrlAlreadyExist
	}
	ms.urls[originalUrl] = shortUrl
	ms.hash[shortUrl] = originalUrl
	return nil
}

func (ms *inMemoryStorage) GetOriginalUrl(_ context.Context, shortUrl string) (string, error) {
	ms.rwMu.RLock()
	defer ms.rwMu.RUnlock()
	if originalUrl, ok := ms.hash[shortUrl]; ok {
		return originalUrl, nil
	}
	return "", errors.New("no rows in result set")
}
