package service

import (
	"encoding/json"
	"fmt"
	"hr-program/internal/attendance-service/dto"
	attdto "hr-program/internal/attendance-service/dto"
	usrdto "hr-program/internal/user-service/dto"
	model "hr-program/shared/models/attendance"
	reqmodel "hr-program/shared/models/request"
	userModel "hr-program/shared/models/users"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

// ดึง attendance logs จาก app DB ผ่าน repository
func (s *AttendanceService) GetAttendanceLogs() ([]model.Attendance, error) {
	return s.AppRepo.GetAttendanceLogs()
}

func (s *AttendanceService) GetAttendanceLogsByDateRange(startDate, endDate string) ([]dto.AttendanceLogsExport, error) {
	attendance, err := s.AppRepo.GetAttendanceLogsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(attendance) == 0 {
		return attendance, nil
	}

	// กรอง duplicate scan: ถ้า user เดียวกันมี SJ ห่างจากครั้งก่อนหน้าไม่เกิน 2 นาที ให้ถือว่าเป็นสแกนซ้ำ
	threshold := 2 * time.Minute
	lastSeen := make(map[string]time.Time) // user_no -> last SJ kept
	filtered := make([]dto.AttendanceLogsExport, 0, len(attendance))

	for _, att := range attendance {
		prev, ok := lastSeen[att.UserNo]
		if ok {
			diff := att.SJ.Sub(prev)
			if diff < 0 {
				diff = -diff
			}
			if diff <= threshold {
				// ภายใน 2 นาทีจากครั้งก่อนของ user เดิม => ถือเป็นสแกนซ้ำ ข้าม และพิมพ์ออก log เพื่อช่วยตรวจสอบ
				log.Printf("duplicate scan user=%s prev=%s current=%s diff=%s", att.UserNo, prev.Format(time.RFC3339), att.SJ.Format(time.RFC3339), diff)
				continue
			}
		}

		lastSeen[att.UserNo] = att.SJ
		filtered = append(filtered, att)
	}
	return filtered, nil
}

func (s *AttendanceService) GetAttendanceDailyByDate(date string) ([]dto.AttendanceDailyDate, error) {
	attendance, err := s.AppRepo.GetAttendanceDailyByDate(date)
	if err != nil {
		return nil, err
	}

	employee, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	empMap := make(map[int64]userModel.Users)
	for _, emp := range employee {
		empMap[emp.ID] = emp
	}

	// track user present
	presentSet := make(map[int64]bool)

	var presentList []dto.PresentDaily
	for _, att := range attendance {
		emp, ok := empMap[att.UserID]
		if !ok {
			continue
		}

		var editedScans []dto.EditedScan
		if len(att.EditedScansJSON) > 0 {
			if err := json.Unmarshal(att.EditedScansJSON, &editedScans); err != nil {
				return nil, err
			}
		}

		present := dto.PresentDaily{
			UserID:          emp.ID,
			EmployeeID:      emp.EmployeeID,
			DepartmentID:    emp.DepartmentID,
			FName:           emp.FName,
			LName:           emp.LName,
			FirstIn:         *att.FirstIn,
			LastOut:         *att.LastOut,
			EditedScansJson: editedScans,
		}
		presentList = append(presentList, present)
		presentSet[emp.ID] = true
	}

	var absentList []dto.AbsentDaily
	for _, emp := range employee {
		if !presentSet[emp.ID] {
			
		}
	}

	return presentList, nil
}

// ดึงและ sync attendance logs จาก Cloudtime -> app DB
func (s *AttendanceService) SyncFullLoadAttendance() error {
	minBH, maxBH, err := s.CloudRepo.GetMinMaxBH()
	if err != nil {
		return err
	}

	mid := minBH + (maxBH-minBH)/2

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// worker 1: [minBH .. mid]
		s.syncRangeAttendance(minBH, mid)
	}()

	go func() {
		defer wg.Done()
		// worker 2: [mid+1 .. maxBH] ป้องกันซ้ำกับ mid ของ worker แรก
		if mid+1 <= maxBH {
			s.syncRangeAttendance(mid+1, maxBH)
		}
	}()

	wg.Wait()

	return nil
}

func (s *AttendanceService) syncRangeAttendance(startBH, endBH int64) {
	batchSize := 3000
	// เริ่มจาก startBH-1 เพื่อให้เงื่อนไข bh > lastBH ครอบคลุม record แรกสุด (bh == startBH)
	lastBH := startBH - 1

	for {
		cloudRecords, err := s.CloudRepo.GetBatchByBHRange(lastBH, endBH, batchSize)
		if err != nil {
			log.Println("Fetch attendance error:", err)
			return
		}

		if len(cloudRecords) == 0 {
			break
		}

		var insertData []model.Attendance
		for _, r := range cloudRecords {
			insertData = append(insertData, model.Attendance{
				BH:           r.BH,
				UserSerial:   r.UserSerial,
				UserNo:       r.UserNo,
				UserLName:    r.UserLName,
				DepNo:        r.DepNo,
				UserDep:      r.UserDep,
				UserDepName:  r.UserDepName,
				UserType:     r.UserType,
				UserCard:     r.UserCard,
				SJ:           r.SJ,
				Iden:         r.Iden,
				FX:           r.FX,
				JlzpSerial:   r.JlzpSerial,
				DevSerial:    r.DevSerial,
				MC:           r.MC,
				HealthStatus: r.HealthStatus,
			})
		}

		if err := s.AppRepo.BulkInsert(insertData); err != nil {
			log.Println("Insert attendance error:", err)
			return
		}

		lastBH = cloudRecords[len(cloudRecords)-1].BH
	}
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
	// 1) ดึง logs ทั้งหมด
	attlogs, err := s.AppRepo.GetAttendanceLogs()
	if err != nil {
		return nil, err
	}

	// 2) เตรียม map employee_id -> user_id
	empIDSet := make(map[string]struct{})
	for _, att := range attlogs {
		empIDSet[att.UserNo] = struct{}{}
	}
	empIDs := make([]string, 0, len(empIDSet))
	for id := range empIDSet {
		empIDs = append(empIDs, id)
	}
	userIDMap, err := s.UserRepo.GetUserIDMapByEmployeeIDs(empIDs)
	if err != nil {
		return nil, err
	}

	// 3) group ตาม (user_id, work_date)
	type groupKey struct {
		UserID   int64
		WorkDate time.Time
	}
	group := make(map[groupKey][]model.Attendance)
	for _, att := range attlogs {
		userID, ok := userIDMap[att.UserNo]
		if !ok {
			continue
		}
		workDate := time.Date(att.SJ.Year(), att.SJ.Month(), att.SJ.Day(), 0, 0, 0, 0, att.SJ.Location())
		key := groupKey{UserID: userID, WorkDate: workDate}
		group[key] = append(group[key], att)
	}

	// 4) เดินทีละ group แล้วค่อยไปหากะ / OT / วันหยุด + คำนวณ daily
	now := time.Now()
	result := make([]model.AttendanceDaily, 0, len(group))

	for key, logs := range group {
		// เรียงสแกนตามเวลา
		sort.Slice(logs, func(i, j int) bool {
			return logs[i].SJ.Before(logs[j].SJ)
		})

		firstIn := logs[0].SJ
		lastOut := logs[len(logs)-1].SJ
		rawJSON, _ := json.Marshal(logs)

		daily := model.AttendanceDaily{
			UserID:           key.UserID,
			WorkDate:         key.WorkDate,
			DayType:          "workday",
			AttendanceStatus: "present",
			FirstIn:          &firstIn,
			LastOut:          &lastOut,
			TotalScans:       len(logs),
			RawScansJSON:     rawJSON,
			CalculatedAt:     ptrTime(now),
		}

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

		editedJSON, err := json.Marshal(editableScans)
		if err != nil {
			return nil, err
		}
		daily.EditedScansJSON = editedJSON
		daily.EditVersion = 0

		// 4.1 หา shift สำหรับวันนั้น
		shifts, err := s.ShiftRepo.GetUserShiftByUserIDAndDate(key.UserID, key.WorkDate)
		if err != nil {
			return nil, err
		}
		if len(shifts) == 0 {
			daily.MissingScan = true
			result = append(result, daily)
			continue
		}
		selected := shifts[0]
		shiftID := selected.ShiftID
		daily.ShiftID = &shiftID
		shiftStart := selected.ShiftDetails.StartTime.Format("15:04:05")
		shiftEnd := selected.ShiftDetails.EndTime.Format("15:04:05")
		daily.ShiftStart = &shiftStart
		daily.ShiftEnd = &shiftEnd
		daily.BreakMinutes = selected.ShiftDetails.BreakMinutes

		// 4.2 หา OT / holiday ตาม key.WorkDate ถ้าต้องใช้
		ots, err := s.OTRepo.GetOTDetailByEmployeeCodeAndDate(key.UserID, key.WorkDate.Format("2006-01-02"))
		if err != nil {
			return nil, err
		}
		// holiday, err := s.HolidayRepo.GetHolidayByDate(key.WorkDate.Format("2026-12-31"))

		shiftMapOT, err := s.matchOTToShift(shifts, ots)
		if err != nil {
			return nil, err
		}

		// 4.3 คำนวณเวลาทำงาน + duplicate scan + ot minutes
		if err := s.calculateWorkMinutes(&daily, shifts); err != nil {
			return nil, err
		}
		if err := s.checkDuplicateScans(&daily); err != nil {
			return nil, err
		}
		if err := s.calculateOTminutes(&daily, shiftMapOT); err != nil {
			return nil, err
		}

		result = append(result, daily)
	}

	return result, nil
}

// func สำหรับ match OT กับกะงาน เพื่อกำหนด relation (before/after/overlap) และเตรียมข้อมูล OTDetailsByDate
func (s *AttendanceService) matchOTToShift(shifts []usrdto.UserShiftAndShiftDetails, ots []reqmodel.OTDetail) (*attdto.UserShiftAndOTByDate, error) {
	if len(ots) == 0 || len(shifts) == 0 {
		return nil, nil // ไม่มีข้อมูลกะงาน ให้ข้ามการคำนวณ OT
	}

	// ใช้กะตัวแรกของวันนั้น (API /shift-user-date คืนมารูปแบบนี้)
	shift := shifts[0]
	shiftDet := shift.ShiftDetails
	otDate := ots[0].Date

	// 1) สร้าง full datetime ของเวลาในกะ: เข้างาน / ออกงาน / ออกเบรค / เข้าเบรค
	shiftStart := time.Date(
		otDate.Year(), otDate.Month(), otDate.Day(),
		shiftDet.StartTime.Hour(), shiftDet.StartTime.Minute(), shiftDet.StartTime.Second(),
		0, otDate.Location(),
	)
	shiftEnd := time.Date(
		otDate.Year(), otDate.Month(), otDate.Day(),
		shiftDet.EndTime.Hour(), shiftDet.EndTime.Minute(), shiftDet.EndTime.Second(),
		0, otDate.Location(),
	)
	breakOut := time.Date(
		otDate.Year(), otDate.Month(), otDate.Day(),
		shiftDet.BreakOut.Hour(), shiftDet.BreakOut.Minute(), shiftDet.BreakOut.Second(),
		0, otDate.Location(),
	)
	breakIn := time.Date(
		otDate.Year(), otDate.Month(), otDate.Day(),
		shiftDet.BreakIn.Hour(), shiftDet.BreakIn.Minute(), shiftDet.BreakIn.Second(),
		0, otDate.Location(),
	)

	// ถ้าเป็นกะดึกและเวลาเลิกน้อยกว่าเวลาเริ่ม ให้เลื่อนวันเลิกไปวันถัดไป
	if shiftDet.IsNightShift && shiftEnd.Before(shiftStart) {
		shiftEnd = shiftEnd.Add(24 * time.Hour)
		breakOut = breakOut.Add(24 * time.Hour)
		breakIn = breakIn.Add(24 * time.Hour)
	}

	// 2) เตรียมก้อนผลลัพธ์ (ยังไม่ใส่ OT)
	shiftDetailsDTO := attdto.ShiftDetails{
		ID:           shiftDet.ID,
		ShiftKey:     shiftDet.ShiftKey,
		ShiftCode:    shiftDet.ShiftCode,
		ShiftName:    shiftDet.ShiftName,
		StartTime:    shiftDet.StartTime,
		EndTime:      shiftDet.EndTime,
		Break:        shiftDet.Break,
		BreakOut:     shiftDet.BreakOut,
		BreakIn:      shiftDet.BreakIn,
		BreakMinutes: shiftDet.BreakMinutes,
		IsNightShift: shiftDet.IsNightShift,
		LivingCost:   shiftDet.LivingCost,
	}

	result := &attdto.UserShiftAndOTByDate{
		UserID:          shift.UserID,
		ShiftID:         shift.ShiftID,
		ShiftDetails:    shiftDetailsDTO,
		StartDate:       shift.StartDate,
		EndDate:         shift.EndDate,
		OTDetailsByDate: make([]attdto.OTDetailsByDate, 0, len(ots)),
	}

	// 3) วนทุก OT ของวันนั้น แล้วคำนวณ relation + append เข้าไป
	for _, ot := range ots {
		// 3.1 สร้างช่วงเวลา OT จาก ot.Date + start_ot/stop_ot
		otStartClock, err := time.Parse("15:04:05", ot.StartOT)
		if err != nil {
			return nil, err
		}
		otStopClock, err := time.Parse("15:04:05", ot.StopOT)
		if err != nil {
			return nil, err
		}

		otStart := time.Date(
			otDate.Year(), otDate.Month(), otDate.Day(),
			otStartClock.Hour(), otStartClock.Minute(), otStartClock.Second(),
			0, otDate.Location(),
		)
		otEnd := time.Date(
			otDate.Year(), otDate.Month(), otDate.Day(),
			otStopClock.Hour(), otStopClock.Minute(), otStopClock.Second(),
			0, otDate.Location(),
		)
		if otEnd.Before(otStart) {
			otEnd = otEnd.Add(24 * time.Hour)
		}

		// 3.2 หาว่า OT นี้ before / after / overlap กับกะ
		relation := "overlap"
		if otEnd.Before(otStart) || otEnd.Equal(shiftStart) {
			relation = "before"
		} else if otStart.After(shiftEnd) || otStart.Equal(shiftEnd) {
			relation = "after"
		}

		// 3.3 สร้าง OTDetailsByDate แล้ว append
		otByDate := attdto.OTDetailsByDate{
			OTID:     int64(ot.ID),
			OTDocID:  ot.OTDocID,
			Date:     ot.Date.Format("2006-01-02"), // layout มาตรฐาน
			StartOT:  ot.StartOT,
			StopOT:   ot.StopOT,
			WorkOT:   ot.WorkOT,
			Relation: relation,
		}

		result.OTDetailsByDate = append(result.OTDetailsByDate, otByDate)
	}

	return result, nil
}

func (s *AttendanceService) calculateOTminutes(daily *model.AttendanceDaily, shiftMapOT *attdto.UserShiftAndOTByDate) error {
	if shiftMapOT == nil {
		daily.OTBeforeMinutes = 0
		daily.OTAfterMinutes = 0
		daily.TotalOTMinutes = 0
		daily.TotalWorkMinutes = daily.NormalWorkMinutes
		return nil
	}

	shiftStart := buildClockOnDate(daily.WorkDate, shiftMapOT.ShiftDetails.StartTime)
	shiftEnd := buildClockOnDate(daily.WorkDate, shiftMapOT.ShiftDetails.EndTime)

	if shiftMapOT.ShiftDetails.IsNightShift && shiftEnd.Before(shiftStart) {
		shiftEnd = shiftEnd.Add(24 * time.Hour)
	}

	beforeMinutes := 0
	afterMinutes := 0

	for _, ot := range shiftMapOT.OTDetailsByDate {
		otStart, err := parseClockStringOnDate(daily.WorkDate, ot.StartOT)
		if err != nil {
			return err
		}
		otEnd, err := parseClockStringOnDate(daily.WorkDate, ot.StopOT)
		if err != nil {
			return err
		}

		if otEnd.Before(otStart) {
			otEnd = otEnd.Add(24 * time.Hour)
		}

		if otStart.Before(shiftStart) {
			beforeEnd := minTime(otEnd, shiftStart)
			if beforeEnd.After(otStart) {
				beforeMinutes += int(beforeEnd.Sub(otStart).Minutes())
			}
		}

		if otEnd.After(shiftEnd) {
			afterStart := maxTime(otStart, shiftEnd)
			if otEnd.After(afterStart) {
				afterMinutes += int(otEnd.Sub(afterStart).Minutes())
			}
		}
	}

	daily.OTBeforeMinutes = beforeMinutes
	daily.OTAfterMinutes = afterMinutes
	daily.TotalOTMinutes = beforeMinutes + afterMinutes
	daily.TotalWorkMinutes = daily.NormalWorkMinutes + daily.TotalOTMinutes

	return nil
}

// func คำนวณเวลาทำงาน และเวลาสาย กลับก่อน จาก EditedScansJSON
func (s *AttendanceService) calculateWorkMinutes(daily *model.AttendanceDaily, shifts []usrdto.UserShiftAndShiftDetails) error {
	if len(shifts) == 0 || daily.ShiftStart == nil || daily.ShiftEnd == nil {
		return nil
	}

	shift := shifts[0].ShiftDetails
	shiftStart := buildClockOnDate(daily.WorkDate, shift.StartTime)
	shiftEnd := buildClockOnDate(daily.WorkDate, shift.EndTime)

	if shift.IsNightShift && shiftEnd.Before(shiftStart) {
		shiftEnd = shiftEnd.Add(24 * time.Hour)
	}

	hasBreak := shift.Break
	var breakStart, breakEnd time.Time
	if hasBreak {
		breakStart = buildClockOnDate(daily.WorkDate, shift.BreakOut)
		breakEnd = buildClockOnDate(daily.WorkDate, shift.BreakIn)

		if shift.IsNightShift && breakEnd.Before(breakStart) {
			breakEnd = breakEnd.Add(24 * time.Hour)
		}
	}

	var scans []model.EditableScan
	if err := json.Unmarshal(daily.EditedScansJSON, &scans); err != nil {
		return err
	}
	if len(scans) == 0 {
		return nil
	}

	sort.Slice(scans, func(i, j int) bool {
		return scans[i].ScanTime.Before(scans[j].ScanTime)
	})

	validScans := make([]model.EditableScan, 0, len(scans))
	for _, scan := range scans {
		if scan.Action != "deleted" {
			validScans = append(validScans, scan)
		}
	}
	if len(validScans) == 0 {
		return nil
	}

	totalNormalMinutes := 0
	var currentIn *time.Time

	for _, scan := range validScans {
		switch scan.Type {
		case "in":
			currentIn = &scan.ScanTime

		case "out":
			if currentIn == nil {
				continue
			}

			intervalStart := *currentIn
			intervalEnd := scan.ScanTime
			if intervalEnd.Before(intervalStart) {
				currentIn = nil
				continue
			}

			normalMinutes := overlapMinutes(intervalStart, intervalEnd, shiftStart, shiftEnd)

			if hasBreak && normalMinutes > 0 {
				normalMinutes -= overlapMinutes(intervalStart, intervalEnd, breakStart, breakEnd)
				if normalMinutes < 0 {
					normalMinutes = 0
				}
			}

			totalNormalMinutes += normalMinutes
			currentIn = nil
		}
	}

	if currentIn != nil {
		daily.MissingScan = true
	}

	maxNormalMinutes := int(shiftEnd.Sub(shiftStart).Minutes())
	if hasBreak {
		maxNormalMinutes -= int(breakEnd.Sub(breakStart).Minutes())
	}
	if maxNormalMinutes < 0 {
		maxNormalMinutes = 0
	}
	if totalNormalMinutes > maxNormalMinutes {
		totalNormalMinutes = maxNormalMinutes
	}

	first := validScans[0]
	last := validScans[len(validScans)-1]

	late := 0
	graceMinutes := 1
	if first.Type == "in" && first.ScanTime.After(shiftStart) {
		diff := int(first.ScanTime.Sub(shiftStart).Minutes())
		if diff > graceMinutes {
			late = diff
		}
	}

	early := 0
	if last.Type == "out" && last.ScanTime.Before(shiftEnd) {
		early = int(shiftEnd.Sub(last.ScanTime).Minutes())
	}

	daily.NormalWorkMinutes = totalNormalMinutes
	daily.LateMinutes = late
	daily.EarlyLeaveMinutes = early
	daily.TotalWorkMinutes = daily.NormalWorkMinutes + daily.TotalOTMinutes

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

func (s *AttendanceService) ExportAttendaceLogsByDateRange(startDate, endDate string) ([]byte, error) {
	attlogs, err := s.AppRepo.GetAttendanceLogsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	builder.WriteString("\uFEFF")

	for index, logItem := range attlogs {
		if index > 0 {
			builder.WriteByte('\n')
		}

		formattedName := strings.Join(strings.Fields(logItem.UserLname), " ")
		formattedIden := strings.TrimSpace(logItem.Iden)
		formattedMC := strings.TrimSpace(logItem.MC)

		builder.WriteString(fmt.Sprintf(
			"%s %s %s %s %s",
			strings.TrimSpace(logItem.UserNo),
			logItem.SJ.Format("02/01/2006 15:04:05"),
			formattedName,
			formattedMC,
			formattedIden,
		))
	}

	return []byte(builder.String()), nil
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

func buildClockOnDate(workDate time.Time, clock time.Time) time.Time {
	return time.Date(
		workDate.Year(),
		workDate.Month(),
		workDate.Day(),
		clock.Hour(),
		clock.Minute(),
		clock.Second(),
		0,
		workDate.Location(),
	)
}

func parseClockStringOnDate(workDate time.Time, value string) (time.Time, error) {
	parsed, err := time.Parse("15:04:05", value)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		workDate.Year(),
		workDate.Month(),
		workDate.Day(),
		parsed.Hour(),
		parsed.Minute(),
		parsed.Second(),
		0,
		workDate.Location(),
	), nil
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func overlapMinutes(aStart, aEnd, bStart, bEnd time.Time) int {
	if aEnd.Before(aStart) || bEnd.Before(bStart) {
		return 0
	}

	start := aStart
	if bStart.After(start) {
		start = bStart
	}

	end := aEnd
	if bEnd.Before(end) {
		end = bEnd
	}

	if !end.After(start) {
		return 0
	}

	return int(end.Sub(start).Minutes())
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
