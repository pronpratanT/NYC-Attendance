package service

import (
	"hr-program/internal/attendance-service/repository"
)

type AttendanceService struct {
	CloudRepo *repository.CloudtimeRepository
	AppRepo   *repository.AttendanceRepository
	UserRepo  repository.UserRepositoryInterface
}

func NewAttendanceService(cloudRepo *repository.CloudtimeRepository, appRepo *repository.AttendanceRepository, userRepo repository.UserRepositoryInterface) *AttendanceService {
	return &AttendanceService{
		CloudRepo: cloudRepo,
		AppRepo:   appRepo,
		UserRepo:  userRepo,
	}
}
