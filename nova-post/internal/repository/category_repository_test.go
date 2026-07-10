package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCategoryRepositoryListCategories(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewCategoryRepository(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`deleted_at` IS NULL ORDER BY path ASC, id ASC")).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"created_at",
			"updated_at",
			"deleted_at",
			"name",
			"parent_id",
			"path",
		}).AddRow(
			1,
			now,
			now,
			nil,
			"Engineering",
			0,
			"/engineering",
		))

	categories, err := repo.ListCategories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) != 1 || categories[0].Name != "Engineering" {
		t.Fatalf("unexpected categories: %+v", categories)
	}
	if categories[0].ID != 1 || categories[0].ParentId != 0 || categories[0].Path != "/engineering" {
		t.Fatalf("unexpected category fields: %+v", categories[0])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
