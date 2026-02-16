package model

import "time"

type Attendance struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	BH           int       `json:"bh"`
	UserSerial   int       `json:"user_serial"`
	UserNo       string    `json:"user_no"`
	UserLName    string    `json:"user_lname"`
	DepNo        string    `json:"dep_no"`
	UserDep      int       `json:"user_dep"`
	UserDepName  string    `json:"user_depname"`
	UserType     int       `json:"user_type"`
	UserCard     string    `json:"user_card"`
	SJ           time.Time `json:"sj"`
	Iden         string    `json:"iden"`
	FX           int       `json:"fx"`
	JlzpSerial   int       `json:"jlzp_serial"`
	DevSerial    string    `json:"dev_serial"`
	MC           string    `json:"mc"`
	HealthStatus int       `json:"health_status"`
	CreatedAt    time.Time `json:"created_at"`
}
