package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxOpenConns    = 0 //unlimited
	maxIdleConns    = 10
	maxIdleTime     = time.Duration(3) * time.Minute
	maxConnLifetime = 0 //unlimited
)

func newDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDB(db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := db.ExecContext(ctx, sqlCreateURLsTable)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)
	db.SetConnMaxLifetime(maxConnLifetime)
	return nil
}

type DatabaseStorage struct {
	db *sql.DB
}

func NewDatabaseStorage(dsn string, timeout time.Duration) (*DatabaseStorage, func() error) {
	db, err := newDB(dsn)
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}

	err = initDB(db, timeout)
	if err != nil {
		log.Fatalf("Can't init database: %v", err)
	}
	return &DatabaseStorage{
		db: db,
	}, db.Close
}

func (dbs *DatabaseStorage) AddURL(ctx context.Context, surl ShortenedURL) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sqlInsertURL, surl.UserID, surl.ShortURL, surl.OrigURL)
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
		_, err := tx.ExecContext(ctx, sqlInsertURL, surl.UserID, surl.ShortURL, surl.OrigURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (dbs *DatabaseStorage) GetOrigURL(ctx context.Context, shortURL string) (string, error) {
	row := dbs.db.QueryRowContext(ctx, sqlGetOrigURL, shortURL)
	var origURL string
	var isDeleted bool
	err := row.Scan(&origURL, &isDeleted)
	if err != nil {
		return "", err
	}
	if isDeleted {
		return "", ErrDeletedURL
	}
	return origURL, nil
}

func (dbs *DatabaseStorage) GetShortURL(ctx context.Context, origURL string) (string, error) {
	row := dbs.db.QueryRowContext(ctx, sqlGetShortURL, origURL)
	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (dbs *DatabaseStorage) Ping(ctx context.Context) error {
	if err := dbs.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (dbs *DatabaseStorage) GetAllUserURLs(ctx context.Context, uid string) ([]ShortenedURL, error) {
	rows, err := dbs.db.QueryContext(ctx, sqlGetAllUserURLs, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var surls []ShortenedURL
	for rows.Next() {
		var surl ShortenedURL
		err = rows.Scan(&surl.UserID, &surl.ShortURL, &surl.OrigURL)
		if err != nil {
			return nil, err
		}
		surls = append(surls, surl)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if surls == nil {
		return nil, ErrNotExist
	}
	return surls, nil
}

func (dbs *DatabaseStorage) DeleteUserURLs(ctx context.Context, dsurls []DelShortURLs) error {
	tx, err := dbs.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, sqlDeleteUserURLs)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, dsurl := range dsurls {
		if _, err := stmt.ExecContext(ctx, dsurl.ShortURL); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
