package service

import (
	deprepository "hr-program/internal/user-service/repository/departments"
	repository "hr-program/internal/user-service/repository/shifts"
	usrrepository "hr-program/internal/user-service/repository/users"
)

type UserService struct {
	CloudtimeUserRepo *usrrepository.CloudtimeUserRepository
	CloudtimeDepRepo  *deprepository.CloudtimeDepartmentsRepository
	AppRepo           *usrrepository.UserRepository
	DepRepo           *deprepository.DepartmentsRepository
	SQLExpressRepo    *repository.SqlExpressShiftRepository
	ShiftRepo         *repository.ShiftsRepository
}

func NewUserService(cloudUserRepo *usrrepository.CloudtimeUserRepository, appRepo *usrrepository.UserRepository, depRepo *deprepository.DepartmentsRepository, cloudDepRepo *deprepository.CloudtimeDepartmentsRepository, sqlExpressRepo *repository.SqlExpressShiftRepository, shiftRepo *repository.ShiftsRepository) *UserService {
	return &UserService{
		CloudtimeUserRepo: cloudUserRepo,
		CloudtimeDepRepo:  cloudDepRepo,
		AppRepo:           appRepo,
		DepRepo:           depRepo,
		SQLExpressRepo:    sqlExpressRepo,
		ShiftRepo:         shiftRepo,
	}
}
