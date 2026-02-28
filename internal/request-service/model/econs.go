package model

import "time"

type OTEcons struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	HRCheck      string    `gorm:"column:hr_check;size:20;not null"`      // pending / approved / rejected
	Sequence     int64     `gorm:"column:sequence;not null;index"`        // SEQUENCE
	Department   string    `gorm:"column:department;size:100;not null"`   // แผนก
	Dep          string    `gorm:"column:dep;not null;index"`             // รหัสแผนก
	ShiftOT      string    `gorm:"column:shift_ot;not null;index"`        // รหัสกะ OT
	TypeOT       string    `gorm:"column:type_ot;size:20;not null"`       // TYPE_OT
	Date         string    `gorm:"column:date;type:date;not null"`        // วันที่ OT
	AB           string    `gorm:"column:ab;size:20;not null"`            // AB
	EmployeeCode string    `gorm:"column:employee_code;size:20;not null"` // รหัสพนักงาน
	StartOT      string    `gorm:"column:start_ot;type:time;not null"`    // เวลาเริ่ม OT
	StopOT       string    `gorm:"column:stop_ot;type:time;not null"`     // เวลาสิ้นสุด OT
	WorkOT       string    `gorm:"column:work_ot;size:255;not null"`      // งานของ OT
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

type HolidayEcons struct {
	ID     int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Date   time.Time `gorm:"column:date;type:date;not null;uniqueIndex"`
	Remark string    `gorm:"column:remark;size:255;"`
	Sunday string    `gorm:"column:sunday;size:10;"` // ใช่/ไม่ใช่
}

func (HolidayEcons) TableName() string {
	return "dbo.HR_HOLIDAY"
}
