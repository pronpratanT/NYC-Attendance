package model

import "time"

type UserShifts struct {
	ID        int64      `gorm:"column:id" json:"id"`
	UserID    int64      `gorm:"column:user_id;not null" json:"user_id"`
	ShiftID   int64      `gorm:"column:shift_id;not null" json:"shift_id"`
	StartDate time.Time  `gorm:"column:start_date;type:date;not null" json:"start_date"`
	EndDate   *time.Time `gorm:"column:end_date;type:date;" json:"end_date"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (UserShifts) TableName() string {
	return "user_shifts"
}
