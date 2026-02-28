package repository

import (
	"hr-program/internal/attendance-service/model"
	"log"

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
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bh"}}, // ใช้ bh เป็น unique key
			DoNothing: true,                          // ถ้า bh ซ้ำ ให้ข้าม record นั้น
		}).
		CreateInBatches(data, len(data)).Error
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
