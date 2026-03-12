package dto

import "time"

type LoginRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}

type UserShiftAndShiftDetails struct {
	UserID       int64        `json:"user_id"`
	ShiftID      int64        `json:"shift_id"`
	ShiftDetails ShiftDetails `json:"shift_details"`
	StartDate    string       `json:"start_date"`
	EndDate      *string      `json:"end_date,omitempty"`
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
