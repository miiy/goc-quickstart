package entity

import "github.com/miiy/goc/db/gorm"

const (
	OwnerTypeUser = "user"

	UploadSceneUnspecified = 0
	UploadSceneAvatar      = 1

	UploadStatusUnspecified = 0
	UploadStatusActive      = 1
	UploadStatusDeleted     = 2
)

type Upload struct {
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
