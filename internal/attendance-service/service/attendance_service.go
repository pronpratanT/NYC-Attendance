package service

import (
	"encoding/json"
	"hr-program/internal/attendance-service/model"
	"sort"
	"time"
)

// ดึง attendance logs จาก app DB ผ่าน repository
func (s *AttendanceService) GetAttendanceLogs() ([]model.Attendance, error) {
	return s.AppRepo.GetAttendanceLogs()
}

// func เรียกใช้การคำนวณ attendance daily และบันทึกลง DB ผ่าน repository
func (s *AttendanceService) GenerateAndSaveAttendanceDaily() error {
	dailies, err := s.AttendanceLogsProcessing()
	if err != nil {
		return err
	}

	// เรียงลำดับข้อมูลตาม work_date เก่าไปใหม่ ก่อนบันทึกลง DB เพื่อให้ข้อมูลใน DB เรียงตาม work_date ด้วย
	sort.Slice(dailies, func(i, j int) bool {
		return dailies[i].WorkDate.Before(dailies[j].WorkDate)
	})
	return s.AppRepo.SaveAttendanceDaily(dailies)
}

// แปลง attendance_logs เป็นกลุ่มต่อคนต่อวัน และเรียงเวลาในแต่ละกลุ่ม
func (s *AttendanceService) AttendanceLogsProcessing() ([]model.AttendanceDaily, error) {
	// ดึง attendance logs จาก app DB ผ่าน repository
	attendanceLogs, err := s.AppRepo.GetAttendanceLogs()
	if err != nil {
		return nil, err
	}

	// collect employeeIDs
	employeeSet := make(map[string]struct{})
	for _, att := range attendanceLogs {
		employeeSet[att.UserNo] = struct{}{}
	}

	employeeIDs := make([]string, 0, len(employeeSet))
	for id := range employeeSet {
		employeeIDs = append(employeeIDs, id)
	}

	// ดึง userID map ทีเดียว
	userMap, err := s.UserRepo.GetUserIDMapByEmployeeIDs(employeeIDs)
	if err != nil {
		return nil, err
	}

	// group user_no + work_date
	type groupKey struct {
		UserID   int64
		WorkDate time.Time
	}

	grouped := make(map[groupKey][]model.Attendance)

	for _, att := range attendanceLogs {
		UserID, ok := userMap[att.UserNo]
		if !ok {
			continue
		}

		workDate := time.Date(
			att.SJ.Year(),
			att.SJ.Month(),
			att.SJ.Day(),
			0, 0, 0, 0,
			att.SJ.Location(),
		)
		key := groupKey{
			UserID:   UserID,
			WorkDate: workDate,
		}
		grouped[key] = append(grouped[key], att)
	}

	now := time.Now()
	result := make([]model.AttendanceDaily, 0, len(grouped))
	for key, logs := range grouped {
		// เรียงเวลา
		sort.Slice(logs, func(i, j int) bool {
			return logs[i].SJ.Before(logs[j].SJ)
		})

		// หาแสกนครั้งแรกและครั้งสุดท้าย
		firstIn := logs[0].SJ
		lastOut := logs[len(logs)-1].SJ

		// แปลง raw logs เป็น JSON
		rawJSON, _ := json.Marshal(logs)

		daily := model.AttendanceDaily{
			UserID:           key.UserID,
			WorkDate:         key.WorkDate,
			DayType:          "workday", // สมมติเป็นวันทำงานก่อน
			AttendanceStatus: "present",

			FirstIn:      &firstIn,
			LastOut:      &lastOut,
			TotalScans:   len(logs),
			RawScansJSON: rawJSON,

			CalculatedAt: ptrTime(now),
		}

		// แปลง attendance logs -> []EditableScan -> EditedScansJSON
		editableScans := make([]model.EditableScan, 0, len(logs))
		for _, l := range logs {
			editableScans = append(editableScans, model.EditableScan{
				ScanTime:  l.SJ,
				Type:      fxToType(l.FX),
				Action:    "added",
				CreatedBy: 0,
				CreatedAt: l.SJ,
			})
		}

		b, err := json.Marshal(editableScans)
		if err != nil {
			return nil, err
		}
		daily.EditedScansJSON = b
		daily.EditVersion = 0

		result = append(result, daily)
	}

	return result, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// helper: แปลง FX จาก attendance_logs เป็น "in"/"out"
func fxToType(fx int) string {
	switch fx {
	case 1:
		return "in"
	case 2:
		return "out"
	default:
		return "in" // หรือค่า default อื่น
	}
}
