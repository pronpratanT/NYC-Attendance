package dto

import "time"

type UserShiftAndOTByDate struct {
	UserID           int64             `json:"user_id"`
	ShiftID          int64             `json:"shift_id"`
	ShiftDetails     ShiftDetails      `json:"shift_details"`
	UserShiftDetails UserShiftDetails  `json:"user_shift_details"`
	StartDate        string            `json:"start_date"`
	EndDate          *string           `json:"end_date,omitempty"`
	OTDetailsByDate  []OTDetailsByDate `json:"ot_details_by_date"`
}

type ShiftDetails struct {
	ID           int64     `json:"id"`
	ShiftKey     int       `json:"shift_key"`
	ShiftCode    string    `json:"shift_code"`
	ShiftName    string    `json:"shift_name"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Break        bool      `json:"break"`
	BreakOut     time.Time `json:"break_out"`
	BreakIn      time.Time `json:"break_in"`
	BreakMinutes int       `json:"break_minutes"`
	IsNightShift bool      `json:"is_night_shift"`
	LivingCost   float64   `json:"living_cost"`
}

type UserShiftDetails struct {
	ID           int64     `json:"id"`
	ShiftKey     int       `json:"shift_key"`
	ShiftCode    string    `json:"shift_code"`
	ShiftName    string    `json:"shift_name"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Break        bool      `json:"break"`
	BreakOut     time.Time `json:"break_out"`
	BreakIn      time.Time `json:"break_in"`
	BreakMinutes int       `json:"break_minutes"`
	IsNightShift bool      `json:"is_night_shift"`
	LivingCost   float64   `json:"living_cost"`
}

type OTDetailsByDate struct {
	OTID     int64  `json:"ot_id"`
	OTDocID  uint   `json:"ot_doc_id"`
	Date     string `json:"date"`
	StartOT  string `json:"start_ot"`
	StopOT   string `json:"stop_ot"`
	WorkOT   string `json:"work_ot"`
	Relation string `json:"relation"` // "before_shift", "after_shift", "holiday", "other"
}

type AttendanceLogsExport struct {
	ID        int64     `json:"id"`
	UserNo    string    `json:"user_no"`
	SJ        time.Time `json:"sj"`
	UserLname string    `json:"user_lname"`
	MC        string    `json:"mc"`
	Iden      string    `json:"iden"`
}

type AttendanceDailyDate struct {
	Date         string         `json:"date"`
	PresentDaily []PresentDaily `json:"present_daily"`
	AbsentDaily  []AbsentDaily  `json:"absent_daily"`
}

type PresentDaily struct {
	UserID          int64        `json:"user_id"`
	EmployeeID      string       `json:"employee_id"`
	DepartmentID    int64        `json:"department_id"`
	FName           string       `json:"f_name"`
	LName           string       `json:"l_name"`
	FirstIn         time.Time    `json:"first_in"`
	LastOut         time.Time    `json:"last_out"`
	EditedScansJson []EditedScan `json:"edited_scans_json"`
}

type EditedScan struct {
	Type      string `json:"type"`
	Action    string `json:"action"`
	ScanTime  string `json:"scan_time"`
	CreatedAt string `json:"created_at"`
	CreatedBy int    `json:"created_by"`
}

type AbsentDaily struct {
	EmployeeID   string `json:"employee_id"`
	DepartmentID int64  `json:"department_id"`
	FName        string `json:"f_name"`
	LName        string `json:"l_name"`
}
