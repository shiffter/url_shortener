package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"url_shortener/config"
	"url_shortener/internal/storage"
)

type postgresStorage struct {
	log    *log.Logger
	dbPool *pgxpool.Pool
}

func NewPostgresStorage(log *log.Logger, cfg *config.Config) storage.Storage {

	pool, err := InitPgxPool(cfg)
	if err != nil {
		log.Fatalf("cannot init postgres %s", err.Error())
	}

	return &postgresStorage{
		log:    log,
		dbPool: pool,
	}
}

func InitPgxPool(c *config.Config) (*pgxpool.Pool, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Pass,
		c.Postgres.DbName,
		c.Postgres.SslMode)
	pool, err := pgxpool.New(context.Background(), connectionUrl)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (ps *postgresStorage) SaveUrlPair(ctx context.Context, shortUrl, originalUrl string) error {
	query := "INSERT INTO url_shortener.urls (short_url, original_url) VALUES ($1, $2)"

	_, err := ps.dbPool.Exec(ctx, query, shortUrl, originalUrl)
	if err != nil {
		ps.log.Error(err)
		return err
	}

	return nil
}

func (ps *postgresStorage) CheckExistingOriginalUrl(ctx context.Context, originalUrl string) (bool, error) {
	query := `SELECT count(*)
				FROM url_shortener.urls
				WHERE url_shortener.urls.original_url = $1
				LIMIT 1
				`

	var rowsCount int
	err := ps.dbPool.QueryRow(ctx, query, originalUrl).Scan(&rowsCount)

	if err != nil {
		return false, err
	}

	return rowsCount > 0, nil
}

func (ps *postgresStorage) CheckExistingShortUrl(ctx context.Context, shortUrl string) (bool, error) {
	query := `SELECT count(*)
				FROM url_shortener.urls
				WHERE url_shortener.urls.short_url = $1
				LIMIT 1
				`

	var rowsCount int
	err := ps.dbPool.QueryRow(ctx, query, shortUrl).Scan(&rowsCount)

	if err != nil {
		return false, err
	}

	return rowsCount > 0, nil
}

func (ps *postgresStorage) GetShortUrl(ctx context.Context, originalUrl string) (string, error) {
	query := `SELECT url_shortener.urls.short_url
				FROM url_shortener.urls
				WHERE url_shortener.urls.original_url = $1
				LIMIT 1
				`

	var shortUrl string
	err := ps.dbPool.QueryRow(ctx, query, originalUrl).Scan(&shortUrl)

	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

func (ps *postgresStorage) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	query := `SELECT url_shortener.urls.original_url
				FROM url_shortener.urls
				WHERE url_shortener.urls.short_url = $1
				LIMIT 1
				`

	var originalUrl string
	err := ps.dbPool.QueryRow(ctx, query, shortUrl).Scan(&originalUrl)

	if err != nil {
		return "", err
	}
	return originalUrl, nil
}
