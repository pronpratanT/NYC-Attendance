package model

import "time"

type Attendance struct {
	ID           int       `gorm:"primaryKey"`
	BH           int64     `gorm:"column:bh;uniqueIndex"`
	UserSerial   int       `gorm:"column:user_serial"`
	UserNo       string    `gorm:"column:user_no;index"`
	UserLName    string    `gorm:"column:user_lname"`
	DepNo        string    `gorm:"column:dep_no"`
	UserDep      int       `gorm:"column:user_dep"`
	UserDepName  string    `gorm:"column:user_depname"`
	UserType     int       `gorm:"column:user_type"`
	UserCard     string    `gorm:"column:user_card"`
	SJ           time.Time `gorm:"column:sj"`
	Iden         string    `gorm:"column:iden"`
	FX           int       `gorm:"column:fx"`
	JlzpSerial   int       `gorm:"column:jlzp_serial"`
	DevSerial    string    `gorm:"column:dev_serial"`
	MC           string    `gorm:"column:mc"`
	HealthStatus int       `gorm:"column:health_status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (Attendance) TableName() string {
	return "attendance_logs"
}
