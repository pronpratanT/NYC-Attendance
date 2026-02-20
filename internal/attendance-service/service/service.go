package service

import (
	"hr-program/internal/attendance-service/repository"
)

type AttendanceService struct {
	CloudRepo *repository.CloudtimeRepository
	AppRepo   *repository.AttendanceRepository
}

func NewAttendanceService(cloudRepo *repository.CloudtimeRepository, appRepo *repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{
		CloudRepo: cloudRepo,
		AppRepo:   appRepo,
	}
}
