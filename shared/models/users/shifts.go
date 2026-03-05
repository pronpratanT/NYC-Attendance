package model

import "time"

type Shifts struct {
	ID           int64      `gorm:"column:id" json:"id"`
	ShiftName    string     `gorm:"column:shift_name;size:50;not null;uniqueIndex:ux_shifts_name" json:"shift_name"`
	StartTime    time.Time  `gorm:"column:start_time;type:time;not null" json:"start_time"`
	EndTime      time.Time  `gorm:"column:end_time;type:time;not null" json:"end_time"`
	Break        bool       `gorm:"column:break;default:false" json:"break"`
	BreakOut     time.Time  `gorm:"column:break_out;type:time" json:"break_out"`
	BreakIn      time.Time  `gorm:"column:break_in;type:time" json:"break_in"`
	BreakMinutes int        `gorm:"column:break_minutes;default:0" json:"break_minutes"`
	IsNightShift bool       `gorm:"column:is_night_shift;default:false" json:"is_night_shift"`
	LivingCost   float64    `gorm:"column:living_cost;type:decimal(10,2);default:0" json:"living_cost"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Shifts) TableName() string {
	return "shifts"
}
