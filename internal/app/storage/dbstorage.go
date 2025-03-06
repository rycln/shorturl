package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func DBInit(databaseDsn string) error {
	var err error
	DB, err = sql.Open("pgx", databaseDsn)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS urls (short_url VARCHAR(7), original_url TEXT)")
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

func (dbs *DatabaseStorage) AddURL(ctx context.Context, shortURL, fullURL string) error {
	_, err := dbs.db.ExecContext(ctx, "INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", shortURL, fullURL)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DatabaseStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	row := dbs.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_url = $1", shortURL)
	var fullURL string
	err := row.Scan(&fullURL)
	if err != nil {
		return "", err
	}
	return fullURL, nil
}
