package server

/*
import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/logger"
	pb "github.com/miiy/goc/component/auth/api/v1"
	"github.com/miiy/goc/component/auth/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

func newMockDb() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true}}, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, err
}

func newMockSrv() pb.AuthServiceServer {
	db, mock, err := newMockDb()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	jwtAuth := jwt.NewJWTAuth(&jwt.Options{
		Secret:    "abcd1234",
		ExpiresIn: 100,
	})
	repo := repository.NewAuthRepository(db)
	logger := logger.NewLogger()
	srv := NewAuthServiceServer(logger, repo, jwtAuth, logger)
	return srv
}

func TestAuthServiceServer_validateRegister(t *testing.T) {
	jsonData := `
{
	"success": [
		{"email": "test@email.com", "username": "test", "password": "test", "password_confirmation": "test"}
	],
	"fail": [
		{"email": "",               "username": "test", "password": "test", "password_confirmation": "test"},
		{"email": "test@email.com", "username": "",     "password": "test", "password_confirmation": "test"},
		{"email": "test@email.com", "username": "test", "password": "",     "password_confirmation": "test"},
		{"email": "test@email.com", "username": "test", "password": "test", "password_confirmation": ""},
		{"email": "test@email.com", "username": "test", "password": "test", "password_confirmation": "test123"}
	]
}
`
	var data map[string][]*pb.RegisterRequest
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range data["success"] {
		err := registerValidate(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, v := range data["fail"] {
		err := registerValidate(v)
		if err == nil {
			log.Fatal(err)
		}
	}
}

func TestAuthServiceServer_Register(t *testing.T) {

	srv, closeFunc := newSrv()
	defer closeFunc()

	ctx := context.Background()

	t.Run("test validate", func(t *testing.T) {
		_, err := srv.Register(ctx, &pb.RegisterRequest{
			Email:                "",
			Username:             "t",
			Password:             "t",
			PasswordConfirmation: "t",
		})
		if err == nil {
			t.Fatal(err)
		}
		t.Log(err)
	})

	mock.ExpectPrepare("SELECT (.+) FROM users")
	cRows := mock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("t", "t").WillReturnRows(cRows)

	mock.ExpectPrepare("INSERT INTO users")
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectPrepare("SELECT (.+) FROM users")
	uRows := mock.NewRows([]string{"id", "username", "password", "status"}).AddRow(1, "t", "", 0)
	mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(1).WillReturnRows(uRows)

	res, err := srv.Register(ctx, &pb.RegisterRequest{
		Email:                "t",
		Username:             "t",
		Password:             "t",
		PasswordConfirmation: "t",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthServiceServer_validateSignIn(t *testing.T) {
	jsonData := `
{
	"success": [
		{"username": "test", "password": "test"}
	],
	"fail": [
		{"username": "", "password": ""}
		{"username": "test", "password": ""}
		{"username": "test", "password": "test123"}
	]
}
`
	var data map[string][]*pb.LoginRequest
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatal(err)
	}
}
*/
