package dto

type OTExport struct {
	EmployeeCode string  `json:"employee_code"`
	Date         string  `json:"date"`
	ShiftID      string  `json:"shift_id"`
	TypeOT       string  `json:"type_ot"`  // "before_shift", "after_shift", "holiday", "other"
	TypeOTs      string  `json:"type_ots"` // "before_shift", "after_shift", "holiday", "other"
	Approve      int     `json:"approve"`  // "pending", "approved", "rejected"
	Hours        float32 `json:"hours"`
}
