package service

import (
	"hr-program/internal/attendance-service/repository"
	usrRepo "hr-program/internal/attendance-service/repository"
	reqRepo "hr-program/internal/request-service/repository"
	shiftRepo "hr-program/internal/user-service/repository/shifts"
)

type AttendanceService struct {
	CloudRepo *repository.CloudtimeRepository
	AppRepo   *repository.AttendanceRepository
	// AppRepo   repository.AttendanceRepositoryInterface
	UserRepo    usrRepo.UserRepositoryInterface
	ShiftRepo   shiftRepo.ShiftRepositoryInterface
	OTRepo      reqRepo.OTRepositoryInterface
	HolidayRepo reqRepo.HolidayRepositoryInterface
}

func NewAttendanceService(
	cloudRepo *repository.CloudtimeRepository,
	appRepo *repository.AttendanceRepository,
	userRepo usrRepo.UserRepositoryInterface,
	shiftRepo shiftRepo.ShiftRepositoryInterface,
	otRepo reqRepo.OTRepositoryInterface,
	holidayRepo reqRepo.HolidayRepositoryInterface,
) *AttendanceService {
	return &AttendanceService{
		CloudRepo:   cloudRepo,
		AppRepo:     appRepo,
		UserRepo:    userRepo,
		ShiftRepo:   shiftRepo,
		OTRepo:      otRepo,
		HolidayRepo: holidayRepo,
	}
}
