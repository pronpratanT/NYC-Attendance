package model

import "time"

type Users struct {
	ID           int64      `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
	EmployeeID   string     `json:"employee_id" gorm:"column:employee_id"`
	Password     string     `json:"password" gorm:"column:password"`
	DepartmentID int64      `json:"department_id" gorm:"column:department_id"`
	FName        string     `json:"f_name" gorm:"column:f_name"`
	LName        string     `json:"l_name" gorm:"column:l_name"`
	IsActive     bool       `json:"is_active" gorm:"column:is_active"`
	Workday      time.Time  `json:"workday" gorm:"column:workday"`
	BirthDate    *time.Time `json:"birth_date" gorm:"column:birth_date"`
}

func (Users) TableName() string {
	return "users"
}
