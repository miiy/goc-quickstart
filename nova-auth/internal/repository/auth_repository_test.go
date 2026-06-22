package repository

import (
	"context"
	"database/sql"
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
