package model

import "time"

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
