package repository

import (
	"context"

	"github.com/miiy/goc-quickstart/nova-user/internal/entity"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/db/gorm/paginate"
)

type UserRepository interface {
	First(ctx context.Context, id int64, columns ...string) (*entity.User, error)
	FindByIDs(ctx context.Context, ids []int64, columns ...string) ([]*entity.User, error)
	List(ctx context.Context, page, pageSize int64, columns ...string) ([]*entity.User, int64, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, id int64, user *entity.User, columns ...string) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) First(ctx context.Context, id int64, columns ...string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Select(columns).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByIDs(ctx context.Context, ids []int64, columns ...string) ([]*entity.User, error) {
	if len(ids) == 0 {
		return []*entity.User{}, nil
	}

	db := r.db.WithContext(ctx)
	if len(columns) > 0 {
		db = db.Select(columns)
	}

	var users []*entity.User
	err := db.Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) List(ctx context.Context, page, pageSize int64, columns ...string) ([]*entity.User, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.User{})

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	total, err := r.findCount(ctx, db)
	if err != nil {
		return nil, 0, err
	}

	scope, _ := paginate.Paginate(int(page), int(pageSize), paginate.DefaultMaxPageSize, int(total))
	var users []*entity.User
	err = db.Scopes(scope).Order("id DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, id int64, v *entity.User, columns ...string) (int64, error) {
	if len(columns) == 0 {
		columns = []string{"nickname", "avatar", "email", "phone", "status", "updated_at"}
	}
	result := r.db.WithContext(ctx).Select(columns).Where("id = ?", id).Updates(v)
	return result.RowsAffected, result.Error
}

func (r *userRepository) Delete(ctx context.Context, id int64) (int64, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.User{})
	return result.RowsAffected, result.Error
}

func (r *userRepository) findCount(ctx context.Context, db *gorm.DB) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Count(&count).Error
	return count, err
}
