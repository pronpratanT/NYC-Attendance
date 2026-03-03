package model

import "time"

type OTDoc struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Sequence   int64     `gorm:"column:sequence;uniqueIndex;not null" json:"sequence"` // อ้างอิงจาก ot_logs.sequence
	Date       time.Time `gorm:"column:date;type:date;not null" json:"date"`
	ShiftOT    string    `gorm:"column:shift_ot;size:10;not null" json:"shift_ot"`
	Department string    `gorm:"column:department;size:100;not null" json:"department"`
	Dep        string    `gorm:"column:dep;size:20;not null" json:"dep"`
	RequestAP  string    `gorm:"column:request_ap;size:100;not null" json:"request_ap"` // ผู้ขอ OT
	RequestTap time.Time `gorm:"column:request_tap" json:"request_tap"`                 // เวลาที่ขอ OT
	HRCheck    string    `gorm:"column:hr_check;size:20" json:"hr_check"`               // pending / approved / rejected ต่อคน
	Approve    string    `gorm:"column:approve;size:20;not null" json:"approve"`        // request / approve / reject

	ChiefAP    string     `gorm:"column:chief_ap;size:100" json:"chief_ap"`     // หัวหน้าที่อนุมัติ
	ChiefTap   *time.Time `gorm:"column:chief_tap" json:"chief_tap"`            // เวลาที่หัวหน้าตอบรับ
	ManagerAP  string     `gorm:"column:manager_ap;size:100" json:"manager_ap"` // ผู้จัดการที่อนุมัติ
	ManagerTap *time.Time `gorm:"column:manager_tap" json:"manager_tap"`        // เวลาที่ผู้จัดการตอบรับ
	HRAP       string     `gorm:"column:hr_ap;size:100" json:"hr_ap"`           // HR ที่อนุมัติ
	HRTap      *time.Time `gorm:"column:hr_tap" json:"hr_tap"`                  // เวลาที่ HR ตอบรับ

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (OTDoc) TableName() string {
	return "ot_doc"
}
