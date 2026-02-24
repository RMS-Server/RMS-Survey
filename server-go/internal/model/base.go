package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel is embedded by all entities that have standard audit columns.
type BaseModel struct {
	ID        string         `gorm:"primaryKey;size:64" json:"id"`
	CreatedAt time.Time      `gorm:"column:create_at" json:"createAt"`
	UpdatedAt time.Time      `gorm:"column:update_at" json:"updateAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:is_deleted;index" json:"-"`
	CreateBy  string         `gorm:"column:create_by;size:256" json:"createBy"`
	UpdateBy  string         `gorm:"column:update_by;size:256" json:"updateBy"`
}

// BaseModelNoSoftDelete is for tables that have audit columns but no soft delete.
type BaseModelNoSoftDelete struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	CreatedAt time.Time `gorm:"column:create_at" json:"createAt"`
	UpdatedAt time.Time `gorm:"column:update_at" json:"updateAt"`
	CreateBy  string    `gorm:"column:create_by;size:256" json:"createBy"`
	UpdateBy  string    `gorm:"column:update_by;size:256" json:"updateBy"`
}

// BaseModelCreateOnly is for tables that only have create_at and create_by.
type BaseModelCreateOnly struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	CreatedAt time.Time `gorm:"column:create_at" json:"createAt"`
	CreateBy  string    `gorm:"column:create_by;size:256" json:"createBy"`
}
