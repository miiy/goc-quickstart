package repository

import (
	"context"

	"github.com/miiy/goc-quickstart/nova-post/internal/entity"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/db/gorm/model"
	"github.com/miiy/goc/db/gorm/paginate"
)

type PostRepository interface {
	First(ctx context.Context, id int64, columns ...string) (*entity.Post, error)
	List(ctx context.Context, filter *ListFilter, page, pageSize int64, columns ...string) ([]*entity.Post, int64, error)
	Create(ctx context.Context, post *entity.Post) error
	Update(ctx context.Context, id int64, post *entity.Post, columns ...string) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

type ListFilter struct {
	AuthorId   int64
	CategoryId int64
	Tag        string
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) First(ctx context.Context, id int64, columns ...string) (*entity.Post, error) {
	var post entity.Post
	db := r.db.WithContext(ctx)
	if len(columns) > 0 {
		db = db.Select(columns)
	}
	err := db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) List(ctx context.Context, filter *ListFilter, page, pageSize int64, columns ...string) ([]*entity.Post, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.Post{})

	// select columns
	if len(columns) > 0 {
		db = db.Select(columns)
	}

	// filter
	if filter.AuthorId > 0 {
		db = db.Where("author_id = ?", filter.AuthorId)
	}
	if filter.CategoryId > 0 {
		db = db.Where("category_id = ?", filter.CategoryId)
	}
	if filter.Tag != "" {
		db = db.Where("tags LIKE ?", "%"+filter.Tag+"%")
	}

	// count
	total, err := r.findCount(ctx, db)
	if err != nil {
		return nil, 0, err
	}

	// paginate
	scope, _ := paginate.Paginate(int(page), int(pageSize), paginate.DefaultMaxPageSize, int(total))
	var posts []*entity.Post
	err = db.Scopes(scope).Order("id DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) Update(ctx context.Context, id int64, v *entity.Post, columns ...string) (int64, error) {
	if len(columns) == 0 {
		// update all fields except auto-managed
		dbNames, err := model.FieldDBNames(v, model.FieldNameExpectAutoSet)
		if err != nil {
			return 0, err
		}
		columns = dbNames
	}
	result := r.db.WithContext(ctx).Select(columns).Where("id = ?", id).Updates(v)
	return result.RowsAffected, result.Error
}

func (r *postRepository) Delete(ctx context.Context, id int64) (int64, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Post{})
	return result.RowsAffected, result.Error
}

func (r *postRepository) findCount(ctx context.Context, db *gorm.DB) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Count(&count).Error
	return count, err
}
