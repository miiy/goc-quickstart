package entity

import "github.com/miiy/goc/db/gorm"

const (
	OwnerTypeUser = "user"

	FileSceneUnspecified = 0
	FileSceneAvatar      = 1

	FileStatusUnspecified = 0
	FileStatusActive      = 1
	FileStatusDeleted     = 2
)

type File struct {
	gorm.Model
	OwnerID   int64
	OwnerType string
	Scene     int64
	ObjectKey string
	URL       string
	MimeType  string
	Size      int64
	Checksum  string
	Status    int64
	CreatedBy int64
}
