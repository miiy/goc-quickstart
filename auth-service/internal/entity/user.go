package entity

import (
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/db/gorm"
)

const (
	UserColumnUsername = "username"
	UserColumnEmail    = "email"
	UserColumnPhone    = "phone"
)

const (
	UserStatusActive   = 1
	UserStatusDisabled = 2
)

type User struct {
	gorm.Model
	Username          string
	Password          string
	Email             string
	EmailVerifiedTime *db.JSONTime
	Phone             string
	Unionid           string
	MpOpenid          string
	MpSessionKey      string
	Status            int64
}
