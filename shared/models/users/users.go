package model

import "time"

type Users struct {
	ID           int64      `gorm:"column:id" json:"id"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	EmployeeID   string     `gorm:"column:employee_id" json:"employee_id"`
	Password     string     `gorm:"column:password" json:"password"`
	DepartmentID int64      `gorm:"column:department_id" json:"department_id"`
	FName        string     `gorm:"column:f_name" json:"f_name"`
	LName        string     `gorm:"column:l_name" json:"l_name"`
	IsActive     bool       `gorm:"column:is_active" json:"is_active"`
	Workday      time.Time  `gorm:"column:workday" json:"workday"`
	BirthDate    *time.Time `gorm:"column:birth_date" json:"birth_date"`
}

func (Users) TableName() string {
	return "users"
}
