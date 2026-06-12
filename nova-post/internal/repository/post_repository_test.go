package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	gormDB, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}}, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return gormDB, mock
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
