package service

import (
	"testing"
	"time"

	usrModel "hr-program/shared/models/users"
)

// Unit test สำหรับ buildUserShiftRecords โดยไม่แตะ DB เลย
func TestBuildUserShiftRecords(t *testing.T) {
	// เตรียม map user: EmpKey -> AppUserID
	users := map[int]int64{
		3589: 439, // EmpKey 3589 map ไปที่ user id 439
	}

	// เตรียม map shift: ShiftKey -> ShiftID
	shiftKeyMap := map[int]int64{
		1029: 109, // ShiftKey 1029 map ไปที่ shift id 109
	}

	// เตรียม master records (TMRESULT)
	tmrDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	masters := []usrModel.SQLExpressMasterKey{
		{
			TmrEmp:  3589, // match users key
			TmrSf:   1029, // match shiftKeyMap key
			TmrDate: tmrDate,
		},
	}

	records := buildUserShiftRecords(users, shiftKeyMap, masters)

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	rec := records[0]
	if rec.UserID != 439 {
		t.Fatalf("expected UserID=439, got %d", rec.UserID)
	}
	if rec.ShiftID != 109 {
		t.Fatalf("expected ShiftID=109, got %d", rec.ShiftID)
	}
	if !rec.StartDate.Equal(time.Now()) {
		t.Fatalf("expected StartDate=%v, got %v", time.Now(), rec.StartDate)
	}
	if rec.EndDate != nil {
		t.Fatalf("expected EndDate to be nil, got %v", rec.EndDate)
	}
	t.Logf("record: UserID=%d ShiftID=%d StartDate=%v EndDate=%v", rec.UserID, rec.ShiftID, rec.StartDate, rec.EndDate)
}
