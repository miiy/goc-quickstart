package entity

import "github.com/miiy/goc/db/gorm"

// Category represents a read-only content category.
type Category struct {
	gorm.Model
	Name     string
	ParentId int64
	Path     string
}
