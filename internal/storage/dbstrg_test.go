package storage

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseStorage_AddURLPair(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func() {
		mock.ExpectClose()

		err = db.Close()
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	}()

	strg := NewDatabaseStorage(db)

	expectedQuery := regexp.QuoteMeta(sqlAddURLPair)

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(expectedQuery).WithArgs(testPair.UID, testPair.Short, testPair.Orig).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := strg.AddURLPair(context.Background(), &testPair)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mock.ExpectBegin()

		err := strg.AddURLPair(ctx, &testPair)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("conflict error", func(t *testing.T) {
		var pgErr = &pgconn.PgError{
			Code: pgerrcode.IntegrityConstraintViolation,
		}

		mock.ExpectBegin()
		mock.ExpectExec(expectedQuery).WithArgs(testPair.UID, testPair.Short, testPair.Orig).WillReturnError(pgErr)

		err := strg.AddURLPair(context.Background(), &testPair)
		assert.ErrorIs(t, err, errConflict)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("tx begin error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errTest)

		err := strg.AddURLPair(context.Background(), &testPair)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabaseStorage_GetURLPairByShort(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func() {
		mock.ExpectClose()

		err = db.Close()
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	}()

	strg := NewDatabaseStorage(db)

	expectedQuery := regexp.QuoteMeta(sqlGetURLPairByShort)

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"user_id", "original_url", "is_deleted"}).AddRow(testPair.UID, testPair.Orig, false)
		mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

		pair, err := strg.GetURLPairByShort(context.Background(), testPair.Short)
		assert.NoError(t, err)
		assert.Equal(t, testPair, *pair)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := strg.GetURLPairByShort(ctx, testPair.Short)
		assert.Error(t, err)
	})

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"user_id", "original_url", "is_deleted"}).AddRow(testPair.UID, testPair.Orig, true)
		mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

		_, err := strg.GetURLPairByShort(context.Background(), testPair.Short)
		assert.ErrorIs(t, err, errDeletedURL)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabaseStorage_AddBatchURLPairs(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func() {
		mock.ExpectClose()

		err = db.Close()
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	}()

	strg := NewDatabaseStorage(db)

	expectedQuery := regexp.QuoteMeta(sqlAddURLPair)

	pairs := []models.URLPair{
		testPair,
	}

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(expectedQuery).WithArgs(testPair.UID, testPair.Short, testPair.Orig).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := strg.AddBatchURLPairs(context.Background(), pairs)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(expectedQuery).WithArgs(testPair.UID, testPair.Short, testPair.Orig).WillReturnError(errTest)

		err := strg.AddBatchURLPairs(context.Background(), pairs)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mock.ExpectBegin()

		err := strg.AddBatchURLPairs(ctx, pairs)
		assert.Error(t, err)
	})

	t.Run("tx begin error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errTest)

		err := strg.AddBatchURLPairs(context.Background(), pairs)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabaseStorage_GetURLPairBatchByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func() {
		mock.ExpectClose()

		err = db.Close()
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	}()

	strg := NewDatabaseStorage(db)

	expectedQuery := regexp.QuoteMeta(sqlGetURLPairBatchByUserID)

	testPairs := []models.URLPair{
		testPair,
	}

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"user_id", "short_url", "original_url"}).AddRow(testPair.UID, testPair.Short, testPair.Orig)
		mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

		pairs, err := strg.GetURLPairBatchByUserID(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testPairs, pairs)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := strg.GetURLPairBatchByUserID(ctx, testUserID)
		assert.Error(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).WillReturnError(errTest)

		_, err := strg.GetURLPairBatchByUserID(context.Background(), testUserID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"user_id", "short_url", "original_url"}))

		_, err := strg.GetURLPairBatchByUserID(context.Background(), testUserID)
		assert.ErrorIs(t, err, errNotExist)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabaseStorage_DeleteRequestedURLs(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func() {
		mock.ExpectClose()

		err = db.Close()
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	}()

	strg := NewDatabaseStorage(db)

	expectedQuery := regexp.QuoteMeta(sqlDeleteRequestedURLs)

	delurls := []*models.DelURLReq{
		&testDelReq,
	}

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectBegin()
		for _, delurl := range delurls {
			mock.ExpectPrepare(expectedQuery).ExpectExec().WithArgs(delurl.Short).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		err := strg.DeleteRequestedURLs(context.Background(), delurls)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("tx begin error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errTest)

		err := strg.DeleteRequestedURLs(context.Background(), delurls)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("prepare error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQuery).WillReturnError(errTest)

		err := strg.DeleteRequestedURLs(context.Background(), delurls)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("exec error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQuery).ExpectExec().WillReturnError(errTest)

		err := strg.DeleteRequestedURLs(context.Background(), delurls)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabaseStorage_Ping(t *testing.T) {
	db, mock, err := sqlmock.New(
		sqlmock.MonitorPingsOption(true),
	)
	require.NoError(t, err)

	defer func() {
		mock.ExpectClose()

		err = db.Close()
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	}()

	strg := NewDatabaseStorage(db)

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(nil)

		err := strg.Ping(context.Background())
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(errTest)

		err := strg.Ping(context.Background())
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
