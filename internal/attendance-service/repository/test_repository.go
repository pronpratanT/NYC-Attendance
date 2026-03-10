package repository

import attmodel "hr-program/shared/models/attendance"

// AttendanceRepositoryInterface ใช้กำหนดสัญญาที่ AttendanceService ต้องการจาก layer repository
// เพื่อให้สามารถสร้าง fake implementation สำหรับ unit test ได้สะดวก
type AttendanceRepositoryInterface interface {
	GetAttendanceLogs() ([]attmodel.Attendance, error)
	SaveAttendanceDaily([]attmodel.AttendanceDaily) error
	BulkInsert([]attmodel.Attendance) error
	GetAttendanceDaily() ([]attmodel.AttendanceDaily, error)
	GetAttendanceDailyByEmployeeID(employeeID int64) ([]attmodel.AttendanceDaily, error)
	GetAttendanceDailyByEmployeeIDAndDate(employeeID int64, date string) ([]attmodel.AttendanceDaily, error)
	GetAttendanceDailyByEmployeeIDAndDateRange(employeeID int64, startDate, endDate string) ([]attmodel.AttendanceDaily, error)
	GetAttendanceDailyByDate(startDate, endDate string) ([]attmodel.AttendanceDaily, error)
}
