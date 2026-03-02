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

type OTDetail struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OTDocID      uint      `gorm:"column:ot_doc_id;index;not null" json:"ot_doc_id"` // FK → ot_doc.id
	EmployeeID   int64     `gorm:"column:employee_id;index" json:"employee_id"`      // user_id ในระบบ (เดิมได้มาจาก employee_code)
	EmployeeCode string    `gorm:"column:employee_code;size:100" json:"employee_code"`
	TypeOT       string    `gorm:"column:type_ot;size:20" json:"type_ot"`
	Date         time.Time `gorm:"column:date;type:date;not null" json:"date"`
	StartOT      time.Time `gorm:"column:start_ot;type:time" json:"start_ot"`
	StopOT       time.Time `gorm:"column:stop_ot;type:time" json:"stop_ot"`
	WorkOT       string    `gorm:"column:work_ot;size:255" json:"work_ot"` // งาน OT โดยรวม

	SourceLogID int64 `gorm:"column:source_log_id;index" json:"source_log_id"` // ot_logs.id
	Sequence    int64 `gorm:"-" json:"sequence"`                               // ใช้เชื่อมกับ OTDoc แต่ไม่เก็บลง DB

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (OTDoc) TableName() string {
	return "ot_doc"
}

func (OTDetail) TableName() string {
	return "ot_detail"
}
