package repository

import (
	"hr-program/internal/attendance-service/model"

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
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "work_date"}}, // ใช้ user_id + work_date เป็น unique key
			DoNothing: true,                                                    // ถ้า user_id + work_date ซ้ำ ให้ข้าม record นั้น
		}).
		CreateInBatches(data, len(data)).Error
}

func (r *AttendanceRepository) GetAttendanceDaily() ([]model.AttendanceDaily, error) {
	var attendance []model.AttendanceDaily
	if err := r.DB.
		Order("work_date ASC").
		Find(&attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}
