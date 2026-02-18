package repository

import (
	"github.com/pronpratanT/leave-system/internal/attendance-service/model"
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
