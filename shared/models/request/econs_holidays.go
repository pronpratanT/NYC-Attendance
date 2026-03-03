package model

import "time"

type HolidayEcons struct {
	ID     int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Date   time.Time `gorm:"column:date;type:date;not null;uniqueIndex" json:"date"`
	Remark string    `gorm:"column:remark;size:255;" json:"remark"`
	Sunday string    `gorm:"column:sunday;size:10;" json:"sunday"` // ใช่/ไม่ใช่
}

func (HolidayEcons) TableName() string {
	return "dbo.HR_HOLIDAY"
}
