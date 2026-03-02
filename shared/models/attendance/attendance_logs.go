package model

import "time"

type Attendance struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	BH           int64     `gorm:"column:bh;uniqueIndex" json:"bh"`
	UserSerial   int       `gorm:"column:user_serial" json:"user_serial"`
	UserNo       string    `gorm:"column:user_no;index" json:"user_no"`
	UserLName    string    `gorm:"column:user_lname" json:"user_lname"`
	DepNo        string    `gorm:"column:dep_no" json:"dep_no"`
	UserDep      int       `gorm:"column:user_dep" json:"user_dep"`
	UserDepName  string    `gorm:"column:user_depname" json:"user_depname"`
	UserType     int       `gorm:"column:user_type" json:"user_type"`
	UserCard     string    `gorm:"column:user_card" json:"user_card"`
	SJ           time.Time `gorm:"column:sj" json:"sj"`
	Iden         string    `gorm:"column:iden" json:"iden"`
	FX           int       `gorm:"column:fx" json:"fx"`
	JlzpSerial   int       `gorm:"column:jlzp_serial" json:"jlzp_serial"`
	DevSerial    string    `gorm:"column:dev_serial" json:"dev_serial"`
	MC           string    `gorm:"column:mc" json:"mc"`
	HealthStatus int       `gorm:"column:health_status" json:"health_status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Attendance) TableName() string {
	return "attendance_logs"
}
