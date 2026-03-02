package model

import "time"

type Departments struct {
	ID        int64      `json:"id" gorm:"primaryKey;column:id"`
	Name      string     `json:"name" gorm:"column:name"`
	DepNo     string     `json:"dep_no" gorm:"column:dep_no"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (Departments) TableName() string {
	return "departments"
}
