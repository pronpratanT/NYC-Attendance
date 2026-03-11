package model

import "time"

type OTDetail struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OTDocID      uint      `gorm:"column:ot_doc_id;index;not null" json:"ot_doc_id"` // FK → ot_doc.id
	EmployeeID   int64     `gorm:"column:employee_id;index" json:"employee_id"`      // user_id ในระบบ (เดิมได้มาจาก employee_code)
	EmployeeCode string    `gorm:"column:employee_code;size:100" json:"employee_code"`
	TypeOT       string    `gorm:"column:type_ot;size:20" json:"type_ot"`
	Date         time.Time `gorm:"column:date;type:date;not null" json:"date"`
	StartOT      string    `gorm:"column:start_ot" json:"start_ot"`
	StopOT       string    `gorm:"column:stop_ot" json:"stop_ot"`
	WorkOT       string    `gorm:"column:work_ot;size:255" json:"work_ot"` // งาน OT โดยรวม

	SourceLogID int64 `gorm:"column:source_log_id;index" json:"source_log_id"` // ot_logs.id
	Sequence    int64 `gorm:"-" json:"sequence"`                               // ใช้เชื่อมกับ OTDoc แต่ไม่เก็บลง DB

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (OTDetail) TableName() string {
	return "ot_detail"
}
