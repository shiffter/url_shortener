package links_usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"url_shortener/internal/storage"
	"url_shortener/internal/usecase"
)

type linksUsecase struct {
	log     *logrus.Logger
	storage storage.Storage
}

func NewLinksUsecase(log *logrus.Logger, storage storage.Storage) usecase.LinksUsecase {
	return &linksUsecase{
		log:     log,
		storage: storage,
	}
}

func (lu *linksUsecase) CreateShortUrl(ctx context.Context, originalUrl string) (string, error) {
	shortUrl := Hasher(originalUrl)
	err := lu.storage.SaveUrlPair(ctx, shortUrl, originalUrl)
	if err != nil {
		lu.log.Error(err)
		return "", err
	}
	return shortUrl, nil
}

func (lu *linksUsecase) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	return lu.storage.GetOriginalUrl(ctx, shortUrl)
}

func Hasher(originalUrl string) string {
	return generateShortURL(originalUrl)
}

func generateShortURL(longURL string) string {
	hash := hashSHA256(longURL)
	encoded := encodeWithCustomAlphabet(hash)
	return encoded
}

func hashSHA256(input string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hasher.Sum(nil)
}

func encodeWithCustomAlphabet(input []byte) string {
	base64Encoded := base64.URLEncoding.EncodeToString(input)
	encoded := make([]byte, 10)
	for i := range encoded { //95 = "_" at ascii
		switch base64Encoded[i] {
		case '+':
			encoded[i] = 95
		case '/':
			encoded[i] = 95
		case '=':
			encoded[i] = 95
		default:
			encoded[i] = base64Encoded[i]
		}
	}
	return string(encoded)
}
