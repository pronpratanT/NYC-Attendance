package migrate

import (
	"log"

	db "hr-program/shared/connection"
	att "hr-program/shared/models/attendance"
	req "hr-program/shared/models/request"
	usr "hr-program/shared/models/users"
)

// AutoMigrate ใช้สำหรับสร้าง/อัปเดตตารางใน app DB จาก shared models
func AutoMigrate() error {
	database := db.ConnectDB()

	if err := database.AutoMigrate(
		&usr.Users{},
		&usr.Shifts{},
		&usr.UserShifts{},
		&usr.UserShiftOverrides{},
		&usr.Departments{},
		&att.Attendance{},
		&att.AttendanceDaily{},
		&req.OTlogs{},
		&req.OTDoc{},
		&req.OTDetail{},
		&req.Holiday{},
	); err != nil {
		log.Println("auto migrate failed:", err)
		return err
	}

	return nil
}
