package model

import "time"

type OTDoc struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Sequence   int64     `gorm:"column:sequence;uniqueIndex;not null"` // อ้างอิงจาก ot_logs.sequence
	Date       time.Time `gorm:"column:date;type:date;not null"`
	ShiftOT    string    `gorm:"column:shift_ot;size:10;not null"`
	Department string    `gorm:"column:department;size:100;not null"`
	Dep        string    `gorm:"column:dep;size:20;not null"`
	RequestAP  string    `gorm:"column:request_ap;size:100;not null"` // ผู้ขอ OT
	RequestTap time.Time `gorm:"column:request_tap"`                  // เวลาที่ขอ OT
	HRCheck    string    `gorm:"column:hr_check;size:20"`             // pending / approved / rejected ต่อคน
	Approve    string    `gorm:"column:approve;size:20;not null"`     // request / approve / reject

	ChiefAP    string     `gorm:"column:chief_ap;size:100"`   // หัวหน้าที่อนุมัติ
	ChiefTap   *time.Time `gorm:"column:chief_tap"`           // เวลาที่หัวหน้าตอบรับ
	ManagerAP  string     `gorm:"column:manager_ap;size:100"` // ผู้จัดการที่อนุมัติ
	ManagerTap *time.Time `gorm:"column:manager_tap"`         // เวลาที่ผู้จัดการตอบรับ
	HRAP       string     `gorm:"column:hr_ap;size:100"`      // HR ที่อนุมัติ
	HRTap      *time.Time `gorm:"column:hr_tap"`              // เวลาที่ HR ตอบรับ

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type OTDetail struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement"`
	OTDocID      uint      `gorm:"column:ot_doc_id;index;not null"` // FK → ot_doc.id
	EmployeeID   int64     `gorm:"column:employee_id;index"`        // user_id ในระบบ (เดิมได้มาจาก employee_code)
	EmployeeName string    `gorm:"column:employee_name;size:100"`
	TypeOT       string    `gorm:"column:type_ot;size:20"`
	StartOT      time.Time `gorm:"column:start_ot;type:time"`
	StopOT       time.Time `gorm:"column:stop_ot;type:time"`
	WorkOT       string    `gorm:"column:work_ot;size:255"` // งาน OT โดยรวม

	SourceLogID int64 `gorm:"column:source_log_id;index"` // ot_logs.id
	Sequence    int64 `gorm:"-"`                          // ใช้เชื่อมกับ OTDoc แต่ไม่เก็บลง DB

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (OTDoc) TableName() string {
	return "ot_doc"
}

func (OTDetail) TableName() string {
	return "ot_detail"
}
