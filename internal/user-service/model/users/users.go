package model

import "time"

type Users struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	EmployeeID   string     `json:"employee_id"`
	Password     string     `json:"password"`
	DepartmentID int64      `json:"department_id"`
	FName        string     `json:"f_name"`
	LName        string     `json:"l_name"`
	IsActive     bool       `json:"is_active"`
	Workday      time.Time  `json:"workday"`
	BirthDate    *time.Time `json:"birth_date"`
}

func (Users) TableName() string {
	return "users"
}
