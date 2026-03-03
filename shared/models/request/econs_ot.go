package model

import "time"

type OTEcons struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	HRCheck      string    `gorm:"column:hr_check;size:20;not null" json:"hr_check"`           // pending / approved / rejected
	Sequence     int64     `gorm:"column:sequence;not null;index" json:"sequence"`             // SEQUENCE
	Department   string    `gorm:"column:department;size:100;not null" json:"department"`      // แผนก
	Dep          string    `gorm:"column:dep;not null;index" json:"dep"`                       // รหัสแผนก
	ShiftOT      string    `gorm:"column:shift_ot;not null;index" json:"shift_ot"`             // รหัสกะ OT
	TypeOT       string    `gorm:"column:type_ot;size:20;not null" json:"type_ot"`             // TYPE_OT
	Date         string    `gorm:"column:date;type:date;not null" json:"date"`                 // วันที่ OT
	AB           string    `gorm:"column:ab;size:20;not null" json:"ab"`                       // AB
	EmployeeCode string    `gorm:"column:employee_code;size:20;not null" json:"employee_code"` // รหัสพนักงาน
	StartOT      string    `gorm:"column:start_ot;type:time;not null" json:"start_ot"`         // เวลาเริ่ม OT
	StopOT       string    `gorm:"column:stop_ot;type:time;not null" json:"stop_ot"`           // เวลาสิ้นสุด OT
	WorkOT       string    `gorm:"column:work_ot;size:255;not null" json:"work_ot"`            // งานของ OT
	Approve      string    `gorm:"column:approve;size:20;not null" json:"approve"`             // approve
	RequestAP    string    `gorm:"column:request_ap;size:100;not null" json:"request_ap"`      // ผู้ขอ OT
	RequestTap   time.Time `gorm:"column:request_tap;autoCreateTime" json:"request_tap"`       // เวลาที่ขอ OT
	ChiefAP      string    `gorm:"column:chief_ap;size:100;not null" json:"chief_ap"`          // หัวหน้าที่อนุมัติ
	ChiefTap     time.Time `gorm:"column:chief_tap;autoCreateTime" json:"chief_tap"`           // เวลาที่หัวหน้าตอบรับ
	ManagerAP    string    `gorm:"column:manager_ap;size:100;not null" json:"manager_ap"`      // ผู้จัดการที่อนุมัติ
	ManagerTap   time.Time `gorm:"column:manager_tap;autoCreateTime" json:"manager_tap"`       // เวลาที่ผู้จัดการตอบรับ
	HRAP         string    `gorm:"column:hr_ap;size:100;not null" json:"hr_ap"`                // HR ที่อนุมัติ
	HRTap        time.Time `gorm:"column:hr_tap;autoCreateTime" json:"hr_tap"`                 // เวลาที่ HR ตอบรับ
	DeleteName   string    `gorm:"column:delete_name;size:100;not null" json:"delete_name"`    // ผู้ที่ลบข้อมูล
	Deletetime   time.Time `gorm:"column:delete_time;autoCreateTime" json:"delete_time"`       // เวลาที่ลบข้อมูล
	CreateDate   time.Time `gorm:"column:create_date;autoCreateTime" json:"create_date"`       // เวลาที่สร้างข้อมูล
}

func (OTEcons) TableName() string {
	return "dbo.HR_OT"
}
