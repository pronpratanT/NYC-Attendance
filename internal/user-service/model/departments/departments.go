package model

import "time"

type Departments struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	DepNo     string     `json:"dep_no"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Departments) TableName() string {
	return "departments"
}
