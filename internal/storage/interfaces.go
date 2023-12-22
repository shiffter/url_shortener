package storage

import "context"

type Storage interface {
	SaveUrlPair(ctx context.Context, shortUrl, originalUrl string) error
	GetShortUrl(ctx context.Context, originalUrl string) (string, error)
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
	CheckExistingOriginalUrl(ctx context.Context, originalUrl string) (bool, error)
	CheckExistingShortUrl(ctx context.Context, shortUrl string) (bool, error)
}
