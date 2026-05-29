package entity

import (
	"github.com/miiy/goc/db/gorm"
)

const (
	PostStatusUnspecified = 0
	PostStatusDraft       = 1
	PostStatusPublished   = 2
)

type Post struct {
	gorm.Model
	AuthorId   int64
	Title      string
	Content    string
	Status     int64
	Tags       string
	CategoryId int64
}
