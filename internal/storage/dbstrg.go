package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/shorturl/internal/models"
)

type DatabaseStorage struct {
	db *sql.DB
}

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{db: db}
}

func (s *DatabaseStorage) AddURLPair(ctx context.Context, pair *models.URLPair) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, sqlAddURLPair, pair.UID, pair.Short, pair.Orig)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return newErrConflict(ErrConflict)
		}
		return err
	}

	return tx.Commit()
}

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
		return nil, newErrDeletedURL(ErrDeletedURL)
	}

	return &pair, nil
}

func (s *DatabaseStorage) AddBatchURLPairs(ctx context.Context, pairs []models.URLPair) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, pair := range pairs {
		_, err := tx.ExecContext(ctx, sqlAddURLPair, pair.UID, pair.Short, pair.Orig)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DatabaseStorage) GetURLPairBatchByUserID(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetURLPairBatchByUserID, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pairs []models.URLPair

	for rows.Next() {
		var pair models.URLPair

		err := rows.Scan(&pair.UID, &pair.Short, &pair.Orig)
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
		return nil, newErrNotExist(ErrNotExist)
	}

	return pairs, nil
}

func (s *DatabaseStorage) DeleteRequestedURLs(ctx context.Context, delurls []*models.DelURLReq) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, sqlDeleteRequestedURLs)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, durl := range delurls {
		if _, err := stmt.ExecContext(ctx, durl.Short); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DatabaseStorage) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *DatabaseStorage) Close() {
	s.db.Close()
}
