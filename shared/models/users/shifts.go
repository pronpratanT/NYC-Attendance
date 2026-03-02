package model

import "time"

type Shifts struct {
	ID           int64      `gorm:"column:id" json:"id"`
	ShiftName    string     `gorm:"column:shift_name;size:50;not null;uniqueIndex:ux_shifts_name" json:"shift_name"`
	StartTime    time.Time  `gorm:"column:start_time;type:time;not null" json:"start_time"`
	EndTime      time.Time  `gorm:"column:end_time;type:time;not null" json:"end_time"`
	BreakMinutes int        `gorm:"column:break_minutes;default:0" json:"break_minutes"`
	IsNightShift bool       `gorm:"column:is_night_shift;default:false" json:"is_night_shift"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Shifts) TableName() string {
	return "shifts"
}

type UserShift struct {
	ID         int64      `gorm:"column:id" json:"id"`
	UserID     int64      `gorm:"column:user_id;not null" json:"user_id"`
	ShiftID    int64      `gorm:"column:shift_id;not null" json:"shift_id"`
	StartDate  time.Time  `gorm:"column:start_date;type:date;not null" json:"start_date"`
	EndDate    time.Time  `gorm:"column:end_date;type:date;not null" json:"end_date"`
	LivingCost int        `gorm:"column:living_cost;default:0" json:"living_cost"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (UserShift) TableName() string {
	return "user_shift"
}

type UserShiftOverrides struct {
	ID        int64      `gorm:"column:id" json:"id"`
	UserID    int64      `gorm:"column:user_id;not null;uniqueIndex:ux_user_override_date" json:"user_id"`
	WorkDate  time.Time  `gorm:"column:work_date;type:date;not null;uniqueIndex:ux_user_override_date" json:"work_date"`
	ShiftID   int64      `gorm:"column:shift_id;not null" json:"shift_id"`
	Reason    string     `gorm:"column:reason;type:text" json:"reason"`
	CreatedBy int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (UserShiftOverrides) TableName() string {
	return "user_shift_overrides"
}
