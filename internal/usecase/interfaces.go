package usecase

import "context"

type LinksUsecase interface {
	CreateShortUrl(ctx context.Context, origUrl string) (string, error)
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
}
