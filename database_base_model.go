package goappbase

import "time"

// gorm.Model alternative without DeletedAt column (to disable Soft Delete)
// see https://gorm.io/docs/delete.html#Soft-Delete
type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
