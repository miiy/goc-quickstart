package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/db/gorm/mysql"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	gormDB, err := openMockGormDB(db)
	if err != nil {
		t.Fatal(err)
	}
	return gormDB, mock
}

func openMockGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
}

func TestPostRepositoryListBindsTagFilter(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostRepository(db)

	tag := "go%' OR 1=1 --"
	tagArg := "%" + tag + "%"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `posts` WHERE tags LIKE ? AND `posts`.`deleted_at` IS NULL")).
		WithArgs(tagArg).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `posts` WHERE tags LIKE ? AND `posts`.`deleted_at` IS NULL ORDER BY id DESC LIMIT ?")).
		WithArgs(tagArg, 10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"deleted_at",
			"author_id",
			"title",
			"content",
			"status",
			"tags",
			"category_id",
		}).AddRow(
			1,
			now,
			now,
			nil,
			42,
			"Post",
			"Content",
			1,
			`["go"]`,
			7,
		))

	posts, total, err := repo.List(context.Background(), &ListFilter{Tag: tag}, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(posts) != 1 {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(posts))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
