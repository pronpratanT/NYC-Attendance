package model

import (
	"encoding/json"
	"time"
)

type AttendanceDaily struct {
	ID int64 `gorm:"primaryKey;column:id"`

	// Identity
	UserID   int64     `gorm:"column:user_id;not null;uniqueIndex:ux_user_work_date"`
	WorkDate time.Time `gorm:"column:work_date;type:date;not null;uniqueIndex:ux_user_work_date"`

	// Day classification
	DayType          string `gorm:"column:day_type;size:20;not null"`          // workday / weekend / holiday
	AttendanceStatus string `gorm:"column:attendance_status;size:20;not null"` // present / late / absent / leave / missing_scan / holiday

	// Shift snapshot
	ShiftStart   *time.Time `gorm:"column:shift_start"` // TIME
	ShiftEnd     *time.Time `gorm:"column:shift_end"`   // TIME
	BreakMinutes int        `gorm:"column:break_minutes;default:0"`

	// Time result
	FirstIn *time.Time `gorm:"column:first_in"` // TIMESTAMP NULL
	LastOut *time.Time `gorm:"column:last_out"`

	// Work breakdown
	NormalWorkMinutes int             `gorm:"column:normal_work_minutes;default:0"`
	OTBeforeMinutes   int             `gorm:"column:ot_before_minutes;default:0"`
	OTAfterMinutes    int             `gorm:"column:ot_after_minutes;default:0"`
	TotalOTMinutes    int             `gorm:"column:total_ot_minutes;default:0"`
	TotalWorkMinutes  int             `gorm:"column:total_work_minutes;default:0"`
	LateMinutes       int             `gorm:"column:late_minutes;default:0"`
	EarlyLeaveMinutes int             `gorm:"column:early_leave_minutes;default:0"`
	TotalScans        int             `gorm:"column:total_scans;default:0"`
	DuplicateScans    int             `gorm:"column:duplicate_scans;default:0"`
	MissingScan       bool            `gorm:"column:missing_scan;default:false"`
	LeaveType         string          `gorm:"column:leave_type;size:20"`
	LeaveMinutes      int             `gorm:"column:leave_minutes;default:0"`
	IsEdited          bool            `gorm:"column:is_edited;default:false"`
	IsLocked          bool            `gorm:"column:is_locked;default:false"`
	RawScansJSON      json.RawMessage `gorm:"column:raw_scans_json;type:jsonb"`

	// System metadata
	CalculatedAt *time.Time `gorm:"column:calculated_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (AttendanceDaily) TableName() string {
	return "attendance_daily"
}
