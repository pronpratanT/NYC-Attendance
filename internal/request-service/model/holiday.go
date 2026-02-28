package model

import "time"

type Holiday struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Date      time.Time `gorm:"column:date;type:date;not null"     json:"date"`
	Remark    string    `gorm:"column:remark"                       json:"remark"` // อนุญาตให้เป็น NULL
	Sunday    string    `gorm:"column:sunday;size:10"               json:"sunday"` // sunday / public หรือ NULL
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"    json:"created_at"`
}

func (Holiday) TableName() string {
	return "holiday"
}
