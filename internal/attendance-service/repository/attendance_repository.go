package repository

import (
	"github.com/pronpratanT/leave-system/internal/attendance-service/model"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

func (r *AttendanceRepository) BulkInsert(data []model.Attendance) error {
	return r.DB.CreateInBatches(data, len(data)).Error
}
