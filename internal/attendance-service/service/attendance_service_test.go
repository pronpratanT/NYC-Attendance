package service

import (
	"hr-program/internal/user-service/dto"
	attmodel "hr-program/shared/models/attendance"
	usermodel "hr-program/shared/models/users"
	"testing"
	"time"
)

type testAppRepo struct {
	logs []attmodel.Attendance
}

func (t *testAppRepo) GetAttendanceLogs() ([]attmodel.Attendance, error) {
	return t.logs, nil
}

// ให้ตรงกับ AttendanceRepositoryInterface แต่ในเทสต์นี้ยังไม่ใช้
func (t *testAppRepo) SaveAttendanceDaily(d []attmodel.AttendanceDaily) error {
	return nil
}

// ให้ตรงกับ AttendanceRepositoryInterface สำหรับ syncRangeAttendance แต่ในเทสต์นี้ไม่ตรวจ BulkInsert
func (t *testAppRepo) BulkInsert(d []attmodel.Attendance) error {
	return nil
}

// ด้านล่างคือ method อ่านข้อมูลที่ handler ใช้ แต่ unit test นี้ยังไม่ตรวจ
func (t *testAppRepo) GetAttendanceDaily() ([]attmodel.AttendanceDaily, error) {
	return nil, nil
}

func (t *testAppRepo) GetAttendanceDailyByEmployeeID(employeeID int64) ([]attmodel.AttendanceDaily, error) {
	return nil, nil
}

func (t *testAppRepo) GetAttendanceDailyByEmployeeIDAndDateRange(employeeID int64, startDate, endDate string) ([]attmodel.AttendanceDaily, error) {
	return nil, nil
}

func (t *testAppRepo) GetAttendanceDailyByDate(startDate, endDate string) ([]attmodel.AttendanceDaily, error) {
	return nil, nil
}

type testUserRepo struct {
	m map[string]int64
}

func (t *testUserRepo) GetUserIDMapByEmployeeIDs(ids []string) (map[string]int64, error) {
	return t.m, nil
}

type testShiftRepo struct {
	result []dto.UserShiftAndShiftDetails
}

func (t *testShiftRepo) GetUserShiftByUserIDAndDate(userID int64, date time.Time) ([]dto.UserShiftAndShiftDetails, error) {
	return t.result, nil
}

// stub ให้ตรงกับ ShiftRepositoryInterface แต่ในเทสต์นี้ยังไม่ใช้
func (t *testShiftRepo) GetUserShiftByUserIDs(userIDs []int64) ([]usermodel.UserShifts, error) {
	return nil, nil
}

func TestAttendanceLogsProcessing_Basic(t *testing.T) {
	// 1) เตรียม attendance logs จำลอง
	workDate := time.Date(2026, 3, 6, 0, 0, 0, 0, time.Local)
	logs := []attmodel.Attendance{
		{
			UserNo: "EMP001",
			SJ:     time.Date(2026, 3, 6, 8, 5, 0, 0, time.Local),
			FX:     1, // in
		},
		{
			UserNo: "EMP001",
			SJ:     time.Date(2026, 3, 6, 17, 0, 0, 0, time.Local),
			FX:     2, // out
		},
	}

	// 2) เตรียม fake repos
	appRepo := &testAppRepo{logs: logs}
	userRepo := &testUserRepo{
		m: map[string]int64{"EMP001": 100},
	}
	shiftRepo := &testShiftRepo{
		result: []dto.UserShiftAndShiftDetails{
			{
				UserID:  100,
				ShiftID: 1,
				ShiftDetails: dto.ShiftDetails{
					StartTime:    time.Date(1, 1, 1, 8, 0, 0, 0, time.UTC),
					EndTime:      time.Date(1, 1, 1, 17, 0, 0, 0, time.UTC),
					BreakMinutes: 60,
				},
			},
		},
	}

	// 3) ประกอบ service
	svc := &AttendanceService{
		AppRepo:   appRepo,
		UserRepo:  userRepo,
		ShiftRepo: shiftRepo,
	}

	// 4) เรียกฟังก์ชันที่ต้องการเทสต์
	dailies, err := svc.AttendanceLogsProcessing()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 5) ตรวจผล
	if len(dailies) != 1 {
		t.Fatalf("expected 1 daily, got %d", len(dailies))
	}
	d := dailies[0]
	if d.UserID != 100 {
		t.Errorf("expected UserID 100, got %d", d.UserID)
	}
	if !d.WorkDate.Equal(workDate) {
		t.Errorf("unexpected WorkDate: %v", d.WorkDate)
	}
	if d.ShiftStart == nil || *d.ShiftStart != "08:00:00" {
		t.Errorf("unexpected ShiftStart: %v", d.ShiftStart)
	}
	if d.ShiftEnd == nil || *d.ShiftEnd != "17:00:00" {
		t.Errorf("unexpected ShiftEnd: %v", d.ShiftEnd)
	}
	if d.BreakMinutes != 60 {
		t.Errorf("unexpected BreakMinutes: %d", d.BreakMinutes)
	}

	// log ข้อมูลที่ได้ไว้ดู (รวม JSON แบบอ่านง่าย)
	t.Logf("daily summary: %+v", d)
	t.Logf("RawScansJSON (string): %s", string(d.RawScansJSON))
	t.Logf("EditedScansJSON (string): %s", string(d.EditedScansJSON))
}
