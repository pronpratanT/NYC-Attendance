package repository

import (
	"hr-program/internal/user-service/dto"
	model "hr-program/shared/models/users"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShiftsRepository struct {
	DB *gorm.DB
}

func NewShiftsRepository(db *gorm.DB) *ShiftsRepository {
	return &ShiftsRepository{DB: db}
}

type ShiftRepositoryInterface interface {
	GetUserShiftByUserIDs(userIDs []int64) ([]model.UserShifts, error)
	GetUserShiftByUserIDAndDate(userID int64, date time.Time) ([]dto.UserShiftAndShiftDetails, error)
}

func (r *ShiftsRepository) GetAllShifts() ([]model.Shifts, error) {
	var shifts []model.Shifts
	err := r.DB.Find(&shifts).Error
	return shifts, err
}

func (r *ShiftsRepository) GetAllUsersShifts() ([]model.UserShifts, error) {
	var usrShifts []model.UserShifts
	err := r.DB.Find(&usrShifts).Error
	return usrShifts, err
}

func (r *ShiftsRepository) GetUserShiftByUserIDs(userIDs []int64) ([]model.UserShifts, error) {
	var usr []model.UserShifts
	err := r.DB.Where("user_id IN ?", userIDs).Find(&usr).Error
	return usr, err
}

func (r *ShiftsRepository) GetUserShiftByUserIDAndDate(userID int64, date time.Time) ([]dto.UserShiftAndShiftDetails, error) {
	var rows []struct {
		UserID       int64      `gorm:"column:user_id"`
		ShiftID      int64      `gorm:"column:shift_id"`
		StartDate    time.Time  `gorm:"column:start_date"`
		EndDate      *time.Time `gorm:"column:end_date"`
		ID           int64      `gorm:"column:id"`
		ShiftKey     int        `gorm:"column:shift_key"`
		ShiftCode    string     `gorm:"column:shift_code"`
		ShiftName    string     `gorm:"column:shift_name"`
		StartTime    time.Time  `gorm:"column:start_time"`
		EndTime      time.Time  `gorm:"column:end_time"`
		Break        bool       `gorm:"column:break"`
		BreakOut     time.Time  `gorm:"column:break_out"`
		BreakIn      time.Time  `gorm:"column:break_in"`
		BreakMinutes int        `gorm:"column:break_minutes"`
		IsNightShift bool       `gorm:"column:is_night_shift"`
		LivingCost   float64    `gorm:"column:living_cost"`
	}
	err := r.DB.Table("user_shifts AS us").
		Joins("JOIN shifts AS s ON us.shift_id = s.id").
		Where(
			"us.user_id = ? AND us.start_date <= ? AND (us.end_date IS NULL OR us.end_date >= ?)",
			userID, date, date,
		).
		Select(`
		us.user_id,
            us.shift_id,
            us.start_date,
            us.end_date,
            s.id,
            s.shift_key,
            s.shift_code,
            s.shift_name,
            s.start_time,
            s.end_time,
            s.break,
            s.break_out,
            s.break_in,
            s.break_minutes,
            s.is_night_shift,
            s.living_cost
		`).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// map ไปเป็น DTO ของคุณ
	var result []dto.UserShiftAndShiftDetails
	for _, rrow := range rows {
		var endDateStr *string
		if rrow.EndDate != nil {
			s := rrow.EndDate.Format("2006-01-02")
			endDateStr = &s
		}
		result = append(result, dto.UserShiftAndShiftDetails{
			UserID:    rrow.UserID,
			ShiftID:   rrow.ShiftID,
			StartDate: rrow.StartDate.Format("2006-01-02"),
			EndDate:   endDateStr,
			ShiftDetails: dto.ShiftDetails{
				ID:           rrow.ID,
				ShiftKey:     rrow.ShiftKey,
				ShiftCode:    rrow.ShiftCode,
				ShiftName:    rrow.ShiftName,
				StartTime:    rrow.StartTime,
				EndTime:      rrow.EndTime,
				Break:        rrow.Break,
				BreakOut:     rrow.BreakOut,
				BreakIn:      rrow.BreakIn,
				BreakMinutes: rrow.BreakMinutes,
				IsNightShift: rrow.IsNightShift,
				LivingCost:   rrow.LivingCost,
			},
		})
	}

	return result, nil
}

func (r *ShiftsRepository) BulkInsertShifts(data []model.Shifts) error {
	return r.DB.
		Clauses(clause.OnConflict{
			// ไม่ระบุ column ให้ Postgres สร้าง "ON CONFLICT DO NOTHING" ทั่วไป
			// ให้จัดการตาม primary key / constraints ที่มีอยู่
			DoNothing: true,
		}).
		CreateInBatches(data, len(data)).Error
}

func (r *ShiftsRepository) BulkInsertUserShifts(data []model.UserShifts) error {
	return r.DB.
		Clauses(clause.OnConflict{
			// Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).
		CreateInBatches(data, len(data)).Error
}

// ใน ShiftsRepository
func (r *ShiftsRepository) GetShiftKeyMap() (map[int]int64, error) {
	var rows []struct {
		ID       int64 `gorm:"column:id"`
		ShiftKey int   `gorm:"column:shift_key"`
	}

	if err := r.DB.
		Model(&model.Shifts{}).
		Select("id, shift_key").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	m := make(map[int]int64)
	for _, row := range rows {
		m[row.ShiftKey] = row.ID
	}
	return m, nil
}
