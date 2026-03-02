package model

import "time"

type Shifts struct {
	ID           int64      `json:"id" gorm:"primaryKey;column:id"`
	ShiftName    string     `json:"shift_name" gorm:"column:shift_name;size:50;not null;uniqueIndex:ux_shifts_name"`
	StartTime    time.Time  `json:"start_time" gorm:"column:start_time;type:time;not null"`
	EndTime      time.Time  `json:"end_time" gorm:"column:end_time;type:time;not null"`
	BreakMinutes int        `json:"break_minutes" gorm:"column:break_minutes;default:0"`
	IsNightShift bool       `json:"is_night_shift" gorm:"column:is_night_shift;default:false"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (Shifts) TableName() string {
	return "shifts"
}

type UserShift struct {
	ID         int64      `json:"id" gorm:"primaryKey;column:id"`
	UserID     int64      `json:"user_id" gorm:"column:user_id;not null"`
	ShiftID    int64      `json:"shift_id" gorm:"column:shift_id;not null"`
	StartDate  time.Time  `json:"start_date" gorm:"column:start_date;type:date;not null"`
	EndDate    time.Time  `json:"end_date" gorm:"column:end_date;type:date;not null"`
	LivingCost int        `json:"living_cost" gorm:"column:living_cost;default:0"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (UserShift) TableName() string {
	return "user_shift"
}

type UserShiftOverrides struct {
	ID        int64      `json:"id" gorm:"primaryKey;column:id"`
	UserID    int64      `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:ux_user_override_date"`
	WorkDate  time.Time  `json:"work_date" gorm:"column:work_date;type:date;not null;uniqueIndex:ux_user_override_date"`
	ShiftID   int64      `json:"shift_id" gorm:"column:shift_id;not null"`
	Reason    string     `json:"reason" gorm:"column:reason;type:text"`
	CreatedBy int64      `json:"created_by" gorm:"column:created_by"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (UserShiftOverrides) TableName() string {
	return "user_shift_overrides"
}
