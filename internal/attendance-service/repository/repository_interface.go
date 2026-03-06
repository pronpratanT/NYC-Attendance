package repository

import (
	"time"

	dto "hr-program/internal/user-service/dto"
)

type UserRepositoryInterface interface {
	GetUserIDMapByEmployeeIDs(employeeIDs []string) (map[string]int64, error)
}

type ShiftRepositoryInterface interface {
	// ดึง user_shifts ดิบตามกลุ่ม user IDs
	GetUserShiftByUserIDs(userIDs []int64) ([]dto.UserShiftAndShiftDetails, error)
	// ดึง user_shifts + shift details ตาม user + วันที่
	GetUserShiftByUserIDAndDate(userID int64, date time.Time) ([]dto.UserShiftAndShiftDetails, error)
}
