package repository

import (
	"context"

	"github.com/miiy/goc-quickstart/nova-file/internal/entity"
	"github.com/miiy/goc/db/gorm"
)

type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *entity.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}
