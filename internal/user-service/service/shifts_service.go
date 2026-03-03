package service

import (
	repository "hr-program/internal/user-service/repository/shifts"
)

type ShiftsService struct {
	SQLExpressRepo *repository.SqlExpressShiftRepository
	AppRepo        *repository.ShiftsRepository
}

func NewShiftsService(sqlExpressRepo *repository.SqlExpressShiftRepository, appRepo *repository.ShiftsRepository) *ShiftsService {
	return &ShiftsService{
		SQLExpressRepo: sqlExpressRepo,
		AppRepo:        appRepo,
	}
}
