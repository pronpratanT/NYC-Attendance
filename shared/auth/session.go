package auth

type SessionData struct {
	UserID     int64  `json:"user_id"`
	EmployeeID string `json:"employee_id"`
	JTI        string `json:"jti"`
	IssuedAt   int64  `json:"issued_at"`
	ExpiresAt  int64  `json:"expires_at"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Revoked    bool   `json:"revoked"`
}
