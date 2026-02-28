package repository

import (
	"hr-program/internal/attendance-service/model"
	"log"
)

func (r *AttendanceRepository) GetAttendanceLogs() ([]model.Attendance, error) {
	sqlDB, err := r.DB.DB()
	if err != nil {
		log.Println("Failed to get raw DB connection:", err)
		return nil, err
	}
	rows, err := sqlDB.Query("SELECT id, bh, user_serial, user_no, user_lname, dep_no, user_dep, user_depname, user_type, user_card, sj, iden, fx, jlzp_serial, dev_serial, mc, health_status, created_at FROM attendance_logs ORDER BY bh DESC")
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var attendance []model.Attendance
	for rows.Next() {
		var att model.Attendance
		err := rows.Scan(
			&att.ID,
			&att.BH,
			&att.UserSerial,
			&att.UserNo,
			&att.UserLName,
			&att.DepNo,
			&att.UserDep,
			&att.UserDepName,
			&att.UserType,
			&att.UserCard,
			&att.SJ,
			&att.Iden,
			&att.FX,
			&att.JlzpSerial,
			&att.DevSerial,
			&att.MC,
			&att.HealthStatus,
			&att.CreatedAt,
		)
		if err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}
		attendance = append(attendance, att)
	}
	return attendance, nil
}
