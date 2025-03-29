package storage

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseStorageAddURL(t *testing.T) {
	type want struct {
		wantErr bool
	}

	tests := []struct {
		name     string
		shortURL string
		origURL  string
		want     want
	}{
		{
			name:     "Simple test #1",
			shortURL: "abcdefg",
			origURL:  "https://practicum.yandex.ru/",
			want: want{
				wantErr: false,
			},
		},
		{
			name:     "Simple test #2",
			shortURL: "1234ABC",
			origURL:  "https://ya.ru/",
			want: want{
				wantErr: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			dbs := &DatabaseStorage{
				db: db,
			}

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO urls").WithArgs(testID, test.shortURL, test.origURL).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			surl := NewShortenedURL(testID, test.shortURL, test.origURL)
			err = dbs.AddURL(context.Background(), surl)
			if test.want.wantErr {
				assert.Error(t, err)
			}
			err = mock.ExpectationsWereMet()
			if test.want.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestDatabaseStorageAddBatchURL(t *testing.T) {
	type want struct {
		mustContain map[string]string
		wantErr     bool
	}

	tests := []struct {
		name      string
		shortURLs []string
		origURLs  []string
		want      want
	}{
		{
			name: "Several pairs of data #1",
			shortURLs: []string{
				"abcdefg",
				"1234ABC",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
				"https://ya.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"abcdefg": "https://practicum.yandex.ru/",
					"1234ABC": "https://ya.ru/",
				},
				wantErr: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			dbs := &DatabaseStorage{
				db: db,
			}

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				mock.ExpectBegin()
				for i := range test.shortURLs {
					mock.ExpectExec("INSERT INTO urls").WithArgs(testID, test.shortURLs[i], test.origURLs[i]).WillReturnResult(sqlmock.NewResult(1, 1))
				}
				mock.ExpectCommit()

				var surls = make([]ShortenedURL, len(test.shortURLs))
				for i := range test.shortURLs {
					surl := NewShortenedURL(testID, test.shortURLs[i], test.origURLs[i])
					surls[i] = surl
				}

				err = dbs.AddBatchURL(context.Background(), surls)
				if test.want.wantErr {
					assert.Error(t, err)
				}
				err = mock.ExpectationsWereMet()
				if test.want.wantErr {
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestDatabaseStorageGetOrigURL(t *testing.T) {
	type want struct {
		wantErr bool
	}

	tests := []struct {
		name      string
		shortURL  string
		origURL   string
		isDeleted bool
		want      want
	}{
		{
			name:      "Simple test #1",
			shortURL:  "abcdefg",
			origURL:   "https://practicum.yandex.ru/",
			isDeleted: false,
			want: want{
				wantErr: false,
			},
		},
		{
			name:      "URL was deleted #1",
			shortURL:  "1234ABC",
			origURL:   "https://ya.ru/",
			isDeleted: true,
			want: want{
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			dbs := &DatabaseStorage{
				db: db,
			}

			rows := mock.NewRows([]string{"original_url", "is_deleted"}).AddRow(test.origURL, test.isDeleted)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)

			orig, err := dbs.GetOrigURL(context.Background(), test.shortURL)

			if !test.want.wantErr {
				assert.Equal(t, test.origURL, orig)
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, ErrDeletedURL)
			}
		})
	}
}

func TestDatabaseStorageGetShortURL(t *testing.T) {
	tests := []struct {
		name      string
		shortURL  string
		origURL   string
		isDeleted bool
	}{
		{
			name:     "Simple test #1",
			shortURL: "abcdefg",
			origURL:  "https://practicum.yandex.ru/",
		},
		{
			name:     "Simple test #1",
			shortURL: "1234ABC",
			origURL:  "https://ya.ru/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			dbs := &DatabaseStorage{
				db: db,
			}

			rows := mock.NewRows([]string{"short_url"}).AddRow(test.shortURL)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)

			short, err := dbs.GetShortURL(context.Background(), test.origURL)

			assert.Equal(t, test.shortURL, short)
			assert.NoError(t, err)
		})
	}
}

func TestDatabaseStoragePing(t *testing.T) {
	t.Run("Ping test without error", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		dbs := &DatabaseStorage{
			db: db,
		}
		mock.ExpectPing().WillReturnError(nil)

		err = dbs.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run("Ping test with error", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		dbs := &DatabaseStorage{
			db: db,
		}
		mock.ExpectPing().WillReturnError(errTest)

		err = dbs.Ping(context.Background())
		assert.Error(t, err)
	})
}

func TestDatabaseStorageGetAllUserURLs(t *testing.T) {
	type want struct {
		mustContain map[string]string
		uid         string
		wantErr     bool
	}

	tests := []struct {
		name      string
		dataUID   string
		shortURLs []string
		origURLs  []string
		want      want
	}{
		{
			name:    "Simple test #1",
			dataUID: testID,
			shortURLs: []string{
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"abcdefg": "https://practicum.yandex.ru/",
				},
				uid:     testID,
				wantErr: false,
			},
		},
		{
			name:    "Simple test #2",
			dataUID: testID,
			shortURLs: []string{
				"1234ABC",
				"235DCE",
			},
			origURLs: []string{
				"https://ya.ru/",
				"https://yandex.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"1234ABC": "https://ya.ru/",
					"235DCE":  "https://yandex.ru/",
				},
				uid:     testID,
				wantErr: false,
			},
		},
		{
			name:    "Wrong GetURL request #1",
			dataUID: "testID",
			shortURLs: []string{
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
			},
			want: want{
				uid:     "2",
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			dbs := &DatabaseStorage{
				db: db,
			}

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				rows := mock.NewRows([]string{"user_id", "short_url", "original_url"})
				for i := range test.shortURLs {
					rows.AddRow(test.dataUID, test.shortURLs[i], test.origURLs[i])
				}
				mock.ExpectQuery("SELECT").WillReturnRows(rows)

				surls, err := dbs.GetAllUserURLs(context.Background(), test.want.uid)
				if err != nil {
					if test.want.wantErr {
						assert.Error(t, err)
					}
				}
				if !test.want.wantErr {
					for _, surl := range surls {
						assert.Equal(t, test.want.mustContain[surl.ShortURL], surl.OrigURL)
					}
				}
			}
		})
	}
}

func TestDatabaseStorageDeleteUserURLs(t *testing.T) {
	type want struct {
		wantErr bool
	}

	tests := []struct {
		name      string
		shortURLs []string
		want      want
	}{
		{
			name: "Simple test #1",
			shortURLs: []string{
				"abcdefg",
			},
			want: want{
				wantErr: false,
			},
		},
		{
			name: "Simple test #2",
			shortURLs: []string{
				"1234ABC",
			},
			want: want{
				wantErr: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			dbs := &DatabaseStorage{
				db: db,
			}

			mock.ExpectBegin()
			for i := range test.shortURLs {
				mock.ExpectExec("UPDATE urls").WithArgs(test.shortURLs[i]).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()

			var dsurls = make([]DelShortURLs, len(test.shortURLs))
			for i := range test.shortURLs {
				dsurls[i] = NewDelShortURLs(testID, test.shortURLs[i])
			}
			err = dbs.DeleteUserURLs(context.Background(), dsurls)
			if test.want.wantErr {
				assert.Error(t, err)
			}
			err = mock.ExpectationsWereMet()
			if test.want.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
