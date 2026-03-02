package model

import "time"

type CloudtimeUser struct {
	UserSerial   int64      `json:"user_serial"`
	UserNo       string     `json:"user_no"`
	UserDep      int64      `json:"user_dep"`
	UserLname    string     `json:"user_lname"`
	IsActive     bool       `json:"is_active"`
	UserWorkday  time.Time  `json:"user_workday"`
	UserBirthday *time.Time `json:"user_birthday"`
}

func (CloudtimeUser) TableName() string {
	return "dt_user"
}
