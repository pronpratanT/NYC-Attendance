package repository

import (
	"time"

	dto "hr-program/internal/user-service/dto"
	reqmodel "hr-program/shared/models/request"
	shfmodel "hr-program/shared/models/users"
)

type UserRepositoryInterface interface {
	GetUserIDMapByEmployeeIDs(employeeIDs []string) (map[string]int64, error)
}

type ShiftRepositoryInterface interface {
	// ดึง user_shifts ดิบตามกลุ่ม user IDs
	GetUserShiftByUserIDs(userIDs []int64) ([]dto.UserShiftAndShiftDetails, error)
	GetShiftByID(shiftID int64) ([]shfmodel.Shifts, error)
	// ดึง user_shifts + shift details ตาม user + ช่วงวันที่
	GetUserShiftByUserIDAndDate(userID int64, date time.Time) ([]dto.UserShiftAndShiftDetails, error)
	GetUserShiftByUserIDAndDateRange(userID int64, dateStart, dateEnd time.Time) ([]dto.UserShiftAndShiftDetails, error)
}

type OTRepositoryInterface interface {
	GetOTDetailByEmployeeCodeAndDate(employeeID int64, date string) ([]reqmodel.OTDetail, error)
}

type HolidayRepositoryInterface interface {
	GetHolidays() ([]reqmodel.Holiday, error)
	GetHolidayByDate(date string) ([]reqmodel.Holiday, error)
	GetHolidayByDateRange(startDate, endDate string) ([]reqmodel.Holiday, error)
}
