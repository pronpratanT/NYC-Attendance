package repository

import "gorm.io/gorm"

type AttendanceLogsRepository struct {
	DB *gorm.DB
}

func NewAttendanceLogsRepository(db *gorm.DB) *AttendanceLogsRepository {
	return &AttendanceLogsRepository{DB: db}
}
