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
