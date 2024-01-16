package storage

import "context"

type Storage interface {
	SaveUrlPair(ctx context.Context, shortUrl, originalUrl string) error
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
}
