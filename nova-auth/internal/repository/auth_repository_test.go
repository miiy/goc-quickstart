package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miiy/goc-quickstart/nova-auth/internal/entity"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/db/gorm/mysql"
)

func newMockDb() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := openMockGormDB(db)
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, err
}

func openMockGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
}

func TestMysqlAuthRepository_Create(t *testing.T) {
	db, mock, err := newMockDb()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewAuthRepository(db)
	err = repo.Create(context.Background(), &entity.User{
		Username:          "test",
		Password:          "123456",
		Email:             "test@test.com",
		EmailVerifiedTime: nil,
		Phone:             "",
		Status:            0,
	})
	if err != nil {
		t.Error(err)
	}

}

func TestAuthRepository_FirstByIdentifierReturnsActiveUser(t *testing.T) {
	db, mock, err := newMockDb()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `id`,`username` FROM `users` WHERE (username = ? AND status = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("alice", entity.UserStatusActive, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(42, "alice"))

	repo := NewAuthRepository(db)
	user, err := repo.FirstByIdentifier(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 42 || user.Username != "alice" {
		t.Fatalf("unexpected user: %+v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
