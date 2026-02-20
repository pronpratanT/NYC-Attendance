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

// แปลง attendance_logs เป็นกลุ่มต่อคนต่อวัน และเรียงเวลาในแต่ละกลุ่ม
func (s *AttendanceService) AttendanceLogsProcessing() ([]model.AttendanceDaily, error) {
	// ดึง attendance logs จาก app DB ผ่าน repository
	attendanceLogs, err := s.AppRepo.GetAttendanceLogs()
	if err != nil {
		return nil, err
	}

	type groupKey struct {
		UserNo   string
		WorkDate time.Time
	}

	grouped := make(map[groupKey][]model.Attendance)

	for _, att := range attendanceLogs {
		UserNo := att.UserNo
		workDate := time.Date(
			att.SJ.Year(),
			att.SJ.Month(),
			att.SJ.Day(),
			0, 0, 0, 0,
			att.SJ.Location(),
		)
		key := groupKey{
			UserNo:   UserNo,
			WorkDate: workDate,
		}
		grouped[key] = append(grouped[key], att)
	}

	result := make([]model.AttendanceDaily, 0, len(grouped))
	for key, logs := range grouped {
		// เรียงเวลา
		sort.Slice(logs, func(i, j int) bool {
			return logs[i].SJ.Before(logs[j].SJ)
		})

		// หาแสกนครั้งแรกและครั้งสุดท้าย
		firstIn := logs[0].SJ
		lastOut := logs[len(logs)-1].SJ

		// แปลง ras logs เป็น JSON
		rawJSON, _ := json.Marshal(logs)

		daily := model.AttendanceDaily{
			WorkDate:         key.WorkDate,
			DayType:          "workday", // สมมติเป็นวันทำงานก่อน
			AttendanceStatus: "present",

			FirstIn:      &firstIn,
			LastOut:      &lastOut,
			TotalScans:   len(logs),
			RawScansJSON: rawJSON,

			CalculatedAt: ptrTime(time.Now()),
		}

		result = append(result, daily)
	}

	return result, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
