package storage

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connection pool parameters
const (
	maxOpenConns    = 0 //unlimited
	maxIdleConns    = 10
	maxIdleTime     = time.Duration(3) * time.Minute
	maxConnLifetime = 0 //unlimited
)

// NewDB creates a new database connection pool with configured settings
func NewDB(dsn string) (*sql.DB, error) {
	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(maxOpenConns)
	database.SetMaxIdleConns(maxIdleConns)
	database.SetConnMaxIdleTime(maxIdleTime)
	database.SetConnMaxLifetime(maxConnLifetime)

	return database, nil
}
