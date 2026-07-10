package entity

import (
	"database/sql"

	"github.com/miiy/goc/db/gorm"
)

const (
	PostStatusUnspecified   = 0
	PostStatusDraft         = 1
	PostStatusPublished     = 2
	PostStatusPendingReview = 3
)

type Post struct {
	gorm.Model
	UserId      int64
	Title       string
	Summary     string
	CoverUrl    string
	Content     string
	Status      int64
	Tags        string
	CategoryId  int64
	PublishedAt sql.NullTime
}
