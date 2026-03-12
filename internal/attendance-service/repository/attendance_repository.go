package repository

import (
	"hr-program/internal/attendance-service/dto"
	model "hr-program/shared/models/attendance"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceRepository struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

func (r *AttendanceRepository) BulkInsert(data []model.Attendance) error {
	// Use a safe batch size to avoid Postgres 65535-parameter limit
	const batchSize = 500

	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bh"}}, // ใช้ bh เป็น unique key
			DoNothing: true,                          // ถ้า bh ซ้ำ ให้ข้าม record นั้น
		}).
		CreateInBatches(data, batchSize).Error
}

func (r *AttendanceRepository) SaveAttendanceDaily(data []model.AttendanceDaily) error {
	// ใช้ batch size เล็กลงเพื่อเลี่ยง limit 65535 parameters ของ Postgres
	const batchSize = 500
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "work_date"}}, // ใช้ user_id + work_date เป็น unique key
			DoNothing: true,                                                    // ถ้า user_id + work_date ซ้ำ ให้ข้าม record นั้น
		}).
		CreateInBatches(data, batchSize).Error
}

func (r *AttendanceRepository) GetAttendanceDaily() ([]model.AttendanceDaily, error) {
	var attendance []model.AttendanceDaily
	if err := r.DB.
		Model(&model.AttendanceDaily{}).
		Order("work_date ASC").
		Find(&attendance).Error; err != nil {
		log.Println("Failed to get attendance daily:", err)
		return nil, err
	}
	return attendance, nil
}

func (r *AttendanceRepository) GetAttendanceDailyByEmployeeID(employeeID int64) ([]model.AttendanceDaily, error) {
	var attendance []model.AttendanceDaily
	if err := r.DB.
		Model(&model.AttendanceDaily{}).
		Where("user_id = ?", employeeID).
		Order("work_date ASC").
		Find(&attendance).Error; err != nil {
		log.Println("Failed to get attendance daily:", err)
		return nil, err
	}
	return attendance, nil
}

// GetAttendanceDailyByEmployeeIDAndDate ดึงข้อมูล attendance_daily ตาม user และวันที่เดียว
func (r *AttendanceRepository) GetAttendanceDailyByEmployeeIDAndDate(employeeID int64, date string) ([]model.AttendanceDaily, error) {
	var attendance []model.AttendanceDaily
	if err := r.DB.
		Model(&model.AttendanceDaily{}).
		Where("user_id = ? AND work_date = ?", employeeID, date).
		Order("work_date ASC").
		Find(&attendance).Error; err != nil {
		log.Println("Failed to get attendance daily by employee and date:", err)
		return nil, err
	}
	return attendance, nil
}

func (r *AttendanceRepository) GetAttendanceDailyByEmployeeIDAndDateRange(employeeID int64, startDate, endDate string) ([]model.AttendanceDaily, error) {
	var attendance []model.AttendanceDaily
	if err := r.DB.
		Model(&model.AttendanceDaily{}).
		Where("user_id = ? AND work_date BETWEEN ? AND ?", employeeID, startDate, endDate).
		Order("work_date ASC").
		Find(&attendance).Error; err != nil {
		log.Println("Failed to get attendance daily:", err)
		return nil, err
	}
	return attendance, nil
}

func (r *AttendanceRepository) GetAttendanceDailyByDate(startDate, endDate string) ([]model.AttendanceDaily, error) {
	var attendance []model.AttendanceDaily
	if err := r.DB.
		Model(&model.AttendanceDaily{}).
		Where("work_date BETWEEN ? AND ?", startDate, endDate).
		Order("work_date ASC").
		Find(&attendance).Error; err != nil {
		log.Println("Failed to get attendance daily:", err)
		return nil, err
	}
	return attendance, nil
}

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

func (r *AttendanceRepository) GetAttendanceLogsByDateRange(startDate, endDate string) ([]dto.AttendanceLogsExport, error) {
	const dateLayout = "2006-01-02"

	start, err := time.Parse(dateLayout, startDate)
	if err != nil {
		log.Println("Invalid start_date:", err)
		return nil, err
	}

	end, err := time.Parse(dateLayout, endDate)
	if err != nil {
		log.Println("Invalid end_date:", err)
		return nil, err
	}

	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endExclusive := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location()).AddDate(0, 0, 1)

	var attendance []dto.AttendanceLogsExport
	if err := r.DB.
		Model(&model.Attendance{}).
		Select("id, user_no, user_lname, sj, mc, iden").
		Where("sj >= ? AND sj < ?", startTime, endExclusive).
		Order("sj ASC").
		Scan(&attendance).Error; err != nil {
		log.Println("Failed to get attendance logs by date range:", err)
		return nil, err
	}

	return attendance, nil
}
