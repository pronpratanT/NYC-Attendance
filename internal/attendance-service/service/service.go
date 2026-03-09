package service

import (
	"hr-program/internal/attendance-service/repository"
	usrRepo "hr-program/internal/attendance-service/repository"
	shiftRepo "hr-program/internal/user-service/repository/shifts"
)

type AttendanceService struct {
	CloudRepo *repository.CloudtimeRepository
	AppRepo   repository.AttendanceRepositoryInterface
	UserRepo  usrRepo.UserRepositoryInterface
	ShiftRepo shiftRepo.ShiftRepositoryInterface
}

func NewAttendanceService(cloudRepo *repository.CloudtimeRepository, appRepo repository.AttendanceRepositoryInterface, userRepo usrRepo.UserRepositoryInterface, shiftRepo shiftRepo.ShiftRepositoryInterface) *AttendanceService {
	return &AttendanceService{
		CloudRepo: cloudRepo,
		AppRepo:   appRepo,
		UserRepo:  userRepo,
		ShiftRepo: shiftRepo,
	}
}
