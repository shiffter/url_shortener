package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"url_shortener/config"
	"url_shortener/internal/customError"
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
	connectionUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.Postgres.User,
		c.Postgres.Pass,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.DbName)
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
		alreadyExistErr := &pgconn.PgError{Code: "23505"}
		if errors.As(err, &alreadyExistErr) {
			return customError.UrlAlreadyExist
		}
		return err
	}
	return nil
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
