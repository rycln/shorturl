package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/shorturl/internal/models"
)

// DatabaseStorage is a PostgreSQL implementation of a URL shortener storage.
type DatabaseStorage struct {
	db *sql.DB
}

// NewDatabaseStorage creates a new DatabaseStorage instance.
func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{db: db}
}

// AddURLPair stores a new URL pair in the database.
func (s *DatabaseStorage) AddURLPair(ctx context.Context, pair *models.URLPair) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = fmt.Errorf("%v; rollback failed: %w", err, rollbackErr)
		}
	}()

	_, err = tx.ExecContext(ctx, sqlAddURLPair, pair.UID, pair.Short, pair.Orig)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return newErrConflict(errConflict)
		}
		return err
	}

	return tx.Commit()
}

// GetURLPairByShort retrieves a URL pair by its short URL.
func (s *DatabaseStorage) GetURLPairByShort(ctx context.Context, short models.ShortURL) (*models.URLPair, error) {
	row := s.db.QueryRowContext(ctx, sqlGetURLPairByShort, short)

	var pair = models.URLPair{
		Short: short,
	}
	var isDeleted bool

	err := row.Scan(&pair.UID, &pair.Orig, &isDeleted)
	if err != nil {
		return nil, err
	}

	if isDeleted {
		return nil, newErrDeletedURL(errDeletedURL)
	}

	return &pair, nil
}

// AddBatchURLPairs stores multiple URL pairs in a single transaction.
func (s *DatabaseStorage) AddBatchURLPairs(ctx context.Context, pairs []models.URLPair) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = fmt.Errorf("%v; rollback failed: %w", err, rollbackErr)
		}
	}()

	for _, pair := range pairs {
		_, err := tx.ExecContext(ctx, sqlAddURLPair, pair.UID, pair.Short, pair.Orig)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetURLPairBatchByUserID retrieves all active URL pairs for a user.
func (s *DatabaseStorage) GetURLPairBatchByUserID(ctx context.Context, uid models.UserID) (pairs []models.URLPair, err error) {
	rows, err := s.db.QueryContext(ctx, sqlGetURLPairBatchByUserID, uid)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rowsCloseErr := rows.Close(); rowsCloseErr != nil {
			err = fmt.Errorf("%v; rows close failed: %w", err, rowsCloseErr)
		}
	}()

	for rows.Next() {
		var pair models.URLPair

		err = rows.Scan(&pair.UID, &pair.Short, &pair.Orig)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, pair)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(pairs) == 0 {
		return nil, newErrNotExist(errNotExist)
	}

	return pairs, nil
}

// DeleteRequestedURLs performs batch soft deletion of URLs.
func (s *DatabaseStorage) DeleteRequestedURLs(ctx context.Context, delurls []*models.DelURLReq) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = fmt.Errorf("%v; rollback failed: %w", err, rollbackErr)
		}
	}()

	stmt, err := tx.PrepareContext(ctx, sqlDeleteRequestedURLs)
	if err != nil {
		return err
	}
	defer func() {
		if stmtCloseErr := stmt.Close(); stmtCloseErr != nil {
			err = fmt.Errorf("%v; stmt close failed: %w", err, stmtCloseErr)
		}
	}()

	for _, durl := range delurls {
		if _, err := stmt.ExecContext(ctx, durl.Short); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Ping verifies database connectivity.
func (s *DatabaseStorage) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

// Close releases all database resources.
func (s *DatabaseStorage) Close() error {
	return s.db.Close()
}
