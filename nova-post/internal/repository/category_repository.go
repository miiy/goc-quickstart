package repository

import (
	"context"

	"github.com/miiy/goc-quickstart/nova-post/internal/entity"
	"github.com/miiy/goc/db/gorm"
)

// CategoryRepository reads post categories independently from post persistence.
type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]*entity.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates the repository for read-only category queries.
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// ListCategories returns categories ordered for tree-style display.
func (r *categoryRepository) ListCategories(ctx context.Context) ([]*entity.Category, error) {
	var categories []*entity.Category
	err := r.db.WithContext(ctx).Order("path ASC, id ASC").Find(&categories).Error
	return categories, err
}
