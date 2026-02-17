package model

import "time"

type CloudtimeAttendance struct {
	BH           int64
	UserSerial   int
	UserNo       string
	UserLName    string
	DepNo        string
	UserDep      int
	UserDepName  string
	UserType     int
	UserCard     string
	SJ           time.Time
	Iden         string
	FX           int
	JlzpSerial   int
	DevSerial    string
	MC           string
	HealthStatus int
}

func (CloudtimeAttendance) TableName() string {
	return "atttime"
}
