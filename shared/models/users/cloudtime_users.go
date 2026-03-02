package model

import "time"

type CloudtimeUser struct {
	UserSerial   int64      `gorm:"column:user_serial" json:"user_serial"`
	UserNo       string     `gorm:"column:user_no" json:"user_no"`
	UserDep      int64      `gorm:"column:user_dep" json:"user_dep"`
	UserLname    string     `gorm:"column:user_lname" json:"user_lname"`
	IsActive     bool       `gorm:"column:is_active" json:"is_active"`
	UserWorkday  time.Time  `gorm:"column:user_workday" json:"user_workday"`
	UserBirthday *time.Time `gorm:"column:user_birthday" json:"user_birthday"`
}

func (CloudtimeUser) TableName() string {
	return "dt_user"
}
