package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func DBInitPostgre(databaseDsn string) error {
	var err error
	DB, err = sql.Open("pgx", databaseDsn)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS urls (short_url VARCHAR(7), original_url TEXT UNIQUE)")
	if err != nil {
		return err
	}
	return nil
}

type DatabaseStorage struct {
	db *sql.DB
}

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{
		db: db,
	}
}

func (dbs *DatabaseStorage) Close() error {
	return dbs.db.Close()
}

func (dbs *DatabaseStorage) AddURL(ctx context.Context, surl ShortenedURL) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", surl.ShortURL, surl.OrigURL)
	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrConflict
		}
		return err
	}
	return tx.Commit()
}

func (dbs *DatabaseStorage) AddBatchURL(ctx context.Context, surls []ShortenedURL) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		return err
	}
	for _, surl := range surls {
		_, err := tx.ExecContext(ctx, "INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", surl.ShortURL, surl.OrigURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (dbs *DatabaseStorage) GetOrigURL(ctx context.Context, shortURL string) (string, error) {
	row := dbs.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_url = $1", shortURL)
	var origURL string
	err := row.Scan(&origURL)
	if err != nil {
		return "", err
	}
	return origURL, nil
}

func (dbs *DatabaseStorage) GetShortURL(ctx context.Context, origURL string) (string, error) {
	row := dbs.db.QueryRowContext(ctx, "SELECT short_url FROM urls WHERE original_url = $1", origURL)
	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}
