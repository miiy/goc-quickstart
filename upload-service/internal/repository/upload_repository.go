package repository

import (
	"context"

	"github.com/miiy/goc-quickstart/upload-service/internal/entity"
	"gorm.io/gorm"
)

type UploadRepository interface {
	Create(ctx context.Context, upload *entity.Upload) error
}

type uploadRepository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) UploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Create(ctx context.Context, upload *entity.Upload) error {
	return r.db.WithContext(ctx).Create(upload).Error
}
