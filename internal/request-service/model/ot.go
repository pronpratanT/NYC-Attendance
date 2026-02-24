package model

import "time"

type OTDetail struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	OTDocumentID int64     `gorm:"column:ot_document_id;not null;index"` // SEQUENCE
	TYPE_OT      string    `gorm:"column:type_ot;size:20;not null"`      // TYPE_OT
	DepartmentID int64     `gorm:"column:department_id;not null;index"`
	ShiftID      int64     `gorm:"column:shift_id;not null;index"` // SHIFT_OT
	EmployeeID   int64     `gorm:"column:employee_id;not null;index"`
	OTDate       string    `gorm:"column:ot_date;type:date;not null"`
	OTStart      string    `gorm:"column:ot_start;type:time;not null"`
	OTEnd        string    `gorm:"column:ot_end;type:time;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

type OT struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement"`
	OTDocumentID    int64     `gorm:"column:ot_document_id;not null;index"`       // SEQUENCE
	TYPE_OT         string    `gorm:"column:type_ot;size:20;not null"`            // TYPE_OT
	OTRequesterName string    `gorm:"column:ot_requester_name;size:100;not null"` // ชื่อผู้ขอ OT
	ChieftainName   string    `gorm:"column:chieftain_name;size:100;not null"`    // ชื่อหัวหน้า
	ManagerName     string    `gorm:"column:manager_name;size:100;not null"`      // ชื่อผู้อนุมัติ
	HRApproveName   string    `gorm:"column:hr_approve_name;size:100;not null"`   // ชื่อ HR ที่อนุมัติ
	HRCheckStatus   string    `gorm:"column:hr_check_status;size:20;not null"`    // pending / approved / rejected
	ApproveStatus   string    `gorm:"column:approve_status;size:20;not null"`     // request / manager_approve / hr_approve / delete
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
}

type OTEcons struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	HRCheck      string    `gorm:"column:hr_check;size:20;not null"`      // pending / approved / rejected
	Sequence     int64     `gorm:"column:sequence;not null;index"`        // SEQUENCE
	Department   string    `gorm:"column:department;size:100;not null"`   // แผนก
	Dep          int64     `gorm:"column:dep;not null;index"`             // รหัสแผนก
	ShiftOT      int64     `gorm:"column:shift_ot;not null;index"`        // รหัสกะ OT
	TypeOT       string    `gorm:"column:type_ot;size:20;not null"`       // TYPE_OT
	Date         string    `gorm:"column:date;type:date;not null"`        // วันที่ OT
	AB           string    `gorm:"column:ab;size:20;not null"`            // AB
	EmployeeCode string    `gorm:"column:employee_code;size:20;not null"` // รหัสพนักงาน
	StartOT      string    `gorm:"column:start_ot;type:time;not null"`    // เวลาเริ่ม OT
	StopOT       string    `gorm:"column:stop_ot;type:time;not null"`     // เวลาสิ้นสุด OT
	WorkOT       string    `gorm:"column:work_ot;size:20;not null"`       // งานของ OT
	Approve      string    `gorm:"column:approve;size:20;not null"`       // approve
	RequestAP    string    `gorm:"column:request_ap;size:100;not null"`   // ผู้ขอ OT
	RequestTap   time.Time `gorm:"column:request_tap;autoCreateTime"`     // เวลาที่ขอ OT
	ChiefAP      string    `gorm:"column:chief_ap;size:100;not null"`     // หัวหน้าที่อนุมัติ
	ChiefTap     time.Time `gorm:"column:chief_tap;autoCreateTime"`       // เวลาที่หัวหน้าตอบรับ
	ManagerAP    string    `gorm:"column:manager_ap;size:100;not null"`   // ผู้จัดการที่อนุมัติ
	ManagerTap   time.Time `gorm:"column:manager_tap;autoCreateTime"`     // เวลาที่ผู้จัดการตอบรับ
	HRAP         string    `gorm:"column:hr_ap;size:100;not null"`        // HR ที่อนุมัติ
	HRTap        time.Time `gorm:"column:hr_tap;autoCreateTime"`          // เวลาที่ HR ตอบรับ
	DeleteName   string    `gorm:"column:delete_name;size:100;not null"`  // ผู้ที่ลบข้อมูล
	Deletetime   time.Time `gorm:"column:delete_time;autoCreateTime"`     // เวลาที่ลบข้อมูล
	CreateDate   time.Time `gorm:"column:create_date;autoCreateTime"`     // เวลาที่สร้างข้อมูล
}

func (OTEcons) TableName() string {
	return "dbo.HR_OT"
}
