package in_memory

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"url_shortener/internal/storage"
)

type inMemoryStorage struct {
	log  *log.Logger
	urls map[string]string
	hash map[string]string
	mu   *sync.Mutex
}

func NewMemoryStorage(log *log.Logger) storage.Storage {
	return &inMemoryStorage{
		log:  log,
		urls: make(map[string]string),
		hash: make(map[string]string),
		mu:   &sync.Mutex{},
	}
}

func (ms *inMemoryStorage) SaveUrlPair(_ context.Context, shortUrl, originalUrl string) error {
	ms.mu.Lock()
	ms.urls[originalUrl] = shortUrl
	ms.hash[shortUrl] = originalUrl
	ms.mu.Unlock()
	return nil
}

func (ms *inMemoryStorage) GetShortUrl(_ context.Context, originalUrl string) (string, error) {
	ms.mu.Lock()
	if shortUrl, ok := ms.urls[originalUrl]; ok {
		return shortUrl, nil
	}
	ms.mu.Unlock()
	return "", errors.New("no rows in result set")
}

func (ms *inMemoryStorage) GetOriginalUrl(_ context.Context, shortUrl string) (string, error) {
	ms.mu.Lock()
	if originalUrl, ok := ms.hash[shortUrl]; ok {
		return originalUrl, nil
	}
	ms.mu.Unlock()
	return "", errors.New("no rows in result set")
}

func (ms *inMemoryStorage) CheckExistingOriginalUrl(_ context.Context, originalUrl string) (bool, error) {
	ms.mu.Lock()
	if _, ok := ms.urls[originalUrl]; ok {
		return ok, nil
	}
	ms.mu.Unlock()
	return false, nil
}

func (ms *inMemoryStorage) CheckExistingShortUrl(ctx context.Context, shortUrl string) (bool, error) {
	ms.mu.Lock()
	if _, ok := ms.hash[shortUrl]; ok {
		return ok, nil
	}
	ms.mu.Unlock()
	return false, nil
}
