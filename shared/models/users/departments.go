package model

import "time"

type Departments struct {
	ID        int64      `gorm:"column:id" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	DepNo     string     `gorm:"column:dep_no" json:"dep_no"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Departments) TableName() string {
	return "departments"
}
