package service

import (
	"encoding/json"
	"fmt"
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
			FirstIn:          &firstIn,
			LastOut:          &lastOut,
			TotalScans:       len(logs),
			RawScansJSON:     rawJSON,
			CalculatedAt:     ptrTime(now),
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

		// =========================
		// 🔹 กำหนด Shift ตรงนี้
		// =========================
		// shift = mockup shift 8:00-17:00
		shift := s.getMockShift(key.UserID, key.WorkDate)

		// เก็บเฉพาะเวลาเป็น string เช่น "08:00:00" เตรียมข้อมูลเป็น string postgres แปลง string -> time.Time ให้เอง
		shiftStart := fmt.Sprintf("%02d:%02d:00", shift.StartHour, shift.StartMinute)
		shiftEnd := fmt.Sprintf("%02d:%02d:00", shift.EndHour, shift.EndMinute)

		daily.ShiftStart = &shiftStart
		daily.ShiftEnd = &shiftEnd
		daily.BreakMinutes = shift.BreakMinutes

		// =========================
		// 🔹 เรียก calculate หลังจากมีข้อมูลครบ
		// =========================
		if err := s.calculateWorkMinutes(&daily); err != nil {
			return nil, err
		}

		// เรียกใช้ฟังก์ชันตรวจสอบการแสกนซ้ำ (duplicate scan) โดยดูจาก EditedScansJSON
		if err := s.checkDuplicateScans(&daily); err != nil {
			return nil, err
		}

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
		return "unknown" // หรือค่า default อื่น
	}
}

// func คำนวณเวลาทำงาน และเวลาสาย กลับก่อน จาก EditedScansJSON
func (s *AttendanceService) calculateWorkMinutes(daily *model.AttendanceDaily) error {
	if daily.ShiftStart == nil || daily.ShiftEnd == nil {
		return nil // ไม่มีข้อมูลกะงาน ไม่สามารถคำนวณได้
	}

	// แปลง string เวลา (HH:MM:SS) ให้เป็น time.Time ตามวันที่ของ WorkDate
	shiftStartTime, err := buildShiftDateTime(daily.WorkDate, *daily.ShiftStart)
	if err != nil {
		return err
	}
	shiftEndTime, err := buildShiftDateTime(daily.WorkDate, *daily.ShiftEnd)
	if err != nil {
		return err
	}

	var scans []model.EditableScan
	if err := json.Unmarshal(daily.EditedScansJSON, &scans); err != nil {
		return err
	}

	if len(scans) == 0 {
		return nil // ไม่มีสแกน ไม่สามารถคำนวณได้
	}

	// เรียงสแกนตามเวลา
	sort.Slice(scans, func(i, j int) bool {
		return scans[i].ScanTime.Before(scans[j].ScanTime)
	})

	// 1. คำนวณ Total Work Minutes
	totalMinutes := 0
	var currentIn *time.Time

	for _, scan := range scans {
		if scan.Action == "deleted" {
			continue // ข้ามสแกนที่ถูกลบ
		}

		switch scan.Type {
		case "in":
			currentIn = &scan.ScanTime
		case "out":
			if currentIn != nil {
				// คำนวณเวลาทำงานระหว่าง currentIn กับ scan.ScanTime
				duration := scan.ScanTime.Sub(*currentIn)
				totalMinutes += int(duration.Minutes())
				currentIn = nil
			}
		}
	}

	// ถ้า in ค้าง -> missing scan
	if currentIn != nil {
		daily.MissingScan = true // มีสแกนเข้าแต่ไม่มีสแกนออก
	}

	// จำกัดเวลาทำงานปกติ สูงสุดไม่เกิน 8 ชั่วโมง (480 นาที)
	if totalMinutes > 480 {
		totalMinutes = 480
	}

	daily.TotalWorkMinutes = totalMinutes
	daily.NormalWorkMinutes = totalMinutes

	// 2. คำนวณ Late Minutes มาสาย
	first := scans[0]
	late := 0
	graceMinutes := 1 // กำหนดเวลายืดหยุ่น 1 นาที

	if first.Type == "in" && first.ScanTime.After(*shiftStartTime) {

		diff := int(first.ScanTime.Sub(*shiftStartTime).Minutes())

		if diff > graceMinutes {
			late = diff
		} else {
			late = 0
		}
	}

	daily.LateMinutes = late

	// 3. คำนวณ Early Leave Minutes กลับก่อน
	last := scans[len(scans)-1]
	early := 0
	if last.Type == "out" && last.ScanTime.Before(*shiftEndTime) {
		early = int(shiftEndTime.Sub(last.ScanTime).Minutes())
	}

	daily.EarlyLeaveMinutes = early

	return nil
}

// func ตรวจสอบการแสกนซ้ำ (duplicate scan) โดยดูจาก EditedScansJSON
// ถ้าเจอสแกนที่มีเวลาเดียวกันและประเภทเดียวกัน (in/out) เกิน 1 ครั้ง ให้ถือว่าเป็น duplicate scan
func (s *AttendanceService) checkDuplicateScans(daily *model.AttendanceDaily) error {

	var scans []model.EditableScan
	// แปลง EditedScansJSON เป็น []EditableScan และนำค่าไปใส่ในตัวแปร scans
	if err := json.Unmarshal(daily.EditedScansJSON, &scans); err != nil {
		return err
	}

	if len(scans) == 0 {
		return nil // ไม่มีสแกน ไม่สามารถตรวจสอบได้
	}

	if len(scans)%2 != 0 {
		daily.MissingScan = true // ถ้าจำนวนสแกนเป็นเลขคี่ แสดงว่าขาดคู่ in/out ลืม scan Missing scan = true
	}

	// ตรวจสอบการแสกรซ้ำโดยใช้ เช็ค type ถ้าเจอ type เดียวกัน ต่อกันเกิน 1 ครั้ง ถือว่าเป็น duplicate scan
	var prevType string
	for _, scan := range scans {
		if scan.Action == "deleted" {
			continue // ข้ามสแกนที่ถูกลบ
		}
		// type ก่อนหน้า เปรียบเทียบกับ type ปัจจุบัน
		if prevType == scan.Type {
			diff := scan.ScanTime.Sub(scan.ScanTime)
			if diff > 0 && diff <= time.Minute {
				// ภายใน 1 นาที ถ้าเจอ type เดียวกันซ้ำกัน ให้ถือว่าเป็น duplicate scan
				scan.Action = "deleted" // update action เป็น deleted เพื่อให้ไม่ถูกนับในการคำนวณเวลาทำงาน
				daily.DuplicateScans++  // ถ้าเหมือนกัน duplicate scan เพิ่มขึ้น 1
			}
		}
		prevType = scan.Type
	}

	// updated EditedScansJSON หลังจากตรวจสอบ duplicate scan แล้ว
	update, err := json.Marshal(scans)
	if err != nil {
		return err
	}
	daily.EditedScansJSON = update

	return nil
}

// buildShiftDateTime แปลงเวลาแบบ HH:MM:SS ให้เป็น time.Time โดยใช้วันที่จาก workDate
func buildShiftDateTime(workDate time.Time, t string) (*time.Time, error) {
	parsed, err := time.Parse("15:04:05", t)
	if err != nil {
		return nil, err
	}
	shift := time.Date(
		workDate.Year(),
		workDate.Month(),
		workDate.Day(),
		parsed.Hour(),
		parsed.Minute(),
		parsed.Second(),
		0,
		workDate.Location(),
	)
	return &shift, nil
}

// ------------Mockup shift---------------
type Shift struct {
	StartHour    int
	StartMinute  int
	EndHour      int
	EndMinute    int
	BreakMinutes int
}

func (s *AttendanceService) getMockShift(userID int64, workDate time.Time) Shift {
	return Shift{
		StartHour:    8,
		StartMinute:  0,
		EndHour:      17,
		EndMinute:    0,
		BreakMinutes: 60,
	}
}

// ---------------------------------------
