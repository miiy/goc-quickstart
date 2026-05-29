package entity

import (
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/db/gorm"
)

const (
	UserStatusUnspecified = 0
	UserStatusActive      = 1
	UserStatusDisable     = 2
)

type User struct {
	gorm.Model
	Username          string
	Password          string
	Nickname          string // TODO: add to DB migration
	Avatar            string // TODO: add to DB migration
	Email             string
	EmailVerifiedTime *db.JSONTime
	Phone             string
	Status            int64
}
