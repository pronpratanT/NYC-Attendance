package service

import (
	"fmt"
	model "hr-program/shared/models/users"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// func (s *UserService) GetLatestShifts(limit int) ([]model.SQLExpressUser, error) {
// 	return s.SQLExpressRepo.GetLatestUserShifts(limit)
// }

// func (s *UserService) GetData() ([]model.SQLExpressMasterKey, error) {
// 	return s.SQLExpressRepo.GetLatestMaster()
// }

// func (s *UserService) GetShifts() ([]model.SQLExpressShifts, error) {
// 	return s.SQLExpressRepo.GetLatestShifts(100)
// }

func (s *UserService) GenerateAndSaveShifts() error {
	// ดึงข้อมูล shift จาก TMSHIFT ผ่าน SQLExpressRepo
	shifts, err := s.ProcessShifts()
	if err != nil {
		return err
	}
	// บันทึกข้อมูล shift ลงตาราง Shifts ใน app DB
	return s.ShiftRepo.BulkInsertShifts(shifts)
}

func (s *UserService) ProcessShifts() ([]model.Shifts, error) {
	// ดึงข้อมูล shift จาก TMSHIFT ผ่าน SQLExpressRepo
	expressShifts, err := s.SQLExpressRepo.GetAllShifts()
	if err != nil {
		return nil, err
	}

	var processedShifts []model.Shifts

	for _, es := range expressShifts {
		// แปลงข้อมูลจาก SQLExpressShifts เป็น Shifts
		processedShifts = append(processedShifts, model.Shifts{
			ShiftKey:  es.SFKey,
			ShiftCode: es.SFCode,
			ShiftName: es.SFName,
			StartTime: es.SFInTime,
			EndTime:   es.SFOutTime,
			Break:     es.SFBreak == "Y",
			BreakOut:  es.SFBrkiTime,
			BreakIn:   es.SFBrkoTime,
			BreakMinutes: CalculateBreakMinutes(model.Shifts{
				Break:    es.SFBreak == "Y",
				BreakOut: es.SFBrkiTime,
				BreakIn:  es.SFBrkoTime,
			}),
			IsNightShift: CalculateIsNightShift(model.Shifts{
				StartTime: es.SF1InTime,
				EndTime:   es.SFOutTime,
			}),
			LivingCost: CalculateLivingCost(model.Shifts{
				ShiftName: es.SFName,
			}),
		})
	}
	return processedShifts, nil
}

// BuildUserIDMapByName แม็ปพนักงานจาก SQL Express (EMPFILE) เข้ากับ Users ใน app DB ด้วยชื่อ-สกุล
// โดย normalize คำนำหน้า (นาย/นาง/นางสาว ฯลฯ) ออกก่อน แล้วคืนค่า map[EmpKey]AppUserID
func (s *UserService) BuildUserIDMapByName(expressUsers []model.SQLExpressUser) (map[int]int64, error) {
	// ดึง users จาก app DB
	appUsers, err := s.AppRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// map ชื่อเต็ม (หลัง normalize) -> app user id
	appByName := make(map[string]int64)
	for _, u := range appUsers {
		firstName := normalizeThaiName(u.FName)
		fullName := strings.TrimSpace(firstName + " " + u.LName)
		fullName = normalizeFullName(firstName, u.LName)
		appByName[fullName] = u.ID
	}

	// ผลลัพธ์: map[EmpKey]AppUserID โดยใช้ชื่อเต็ม (หลัง normalize) เป็นตัวเชื่อม
	result := make(map[int]int64)
	for _, eu := range expressUsers {
		fullName := strings.TrimSpace(eu.EmpName + " " + eu.EmpSurnme)
		fullName = normalizeFullName(eu.EmpName, eu.EmpSurnme)
		if id, ok := appByName[fullName]; ok {
			result[eu.EmpKey] = id
		}
	}

	return result, nil
}

func (s *UserService) ProcessUserShifts() error {
	// ดึง user จาก SQL Express แล้วสร้าง map[EmpKey]AppUserID
	expressUsers, err := s.SQLExpressRepo.GetUserBplus()
	if err != nil {
		return err
	}

	users, err := s.BuildUserIDMapByName(expressUsers)
	if err != nil {
		return err
	}

	// map[ShiftKey]ShiftID จากตาราง shifts ใน app DB
	shiftKeyMap, err := s.ShiftRepo.GetShiftKeyMap()
	if err != nil {
		return err
	}

	sqlMaster, err := s.SQLExpressRepo.GetMasterKey()
	if err != nil {
		return err
	}

	records := buildUserShiftRecords(users, shiftKeyMap, sqlMaster)
	if len(records) == 0 {
		return nil
	}

	return s.ShiftRepo.BulkInsertUserShifts(records)
}

// buildUserShiftRecords คือ pure function ที่รับข้อมูลที่เตรียมไว้แล้ว
// (map user, map shift, และ master records) แล้วคืน slice ของ UserShifts
func buildUserShiftRecords(users map[int]int64, shiftKeyMap map[int]int64, masters []model.SQLExpressMasterKey) []model.UserShifts {
	var records []model.UserShifts
	seen := make(map[string]struct{}) // key: userID-shiftID-date

	for _, master := range masters {
		userID, ok := users[master.TmrEmp]
		if !ok {
			continue
		}

		shiftID, ok := shiftKeyMap[master.TmrSf]
		if !ok {
			continue
		}

		now := time.Now()
		dateOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		key := fmt.Sprintf("%d-%d-%s", userID, shiftID, dateOnly.Format("2006-01-02"))
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		records = append(records, model.UserShifts{
			UserID:    userID,
			ShiftID:   shiftID,
			StartDate: dateOnly,
			EndDate:   nil,
		})
	}

	return records
}

// ตัดคำนำหน้าชื่อภาษาไทยพื้นฐานออก เช่น นาย, นาง, นางสาว, น.ส.
func normalizeThaiName(name string) string {
	name = strings.TrimSpace(name)
	prefixes := []string{"นาย", "นางสาว", "น.ส.", "น.ส", "นาง"}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			name = strings.TrimSpace(name[len(p):])
			break
		}
	}
	return name
}

// ตัดเว้นวรรคระหว่างชื่อและสกุลให้เหลือช่องเดียว
func normalizeFullName(first, last string) string {
	name := strings.TrimSpace(first) + " " + strings.TrimSpace(last)
	// บีบ space ซ้ำให้เหลือช่องเดียว
	name = strings.Join(strings.Fields(name), " ")
	return name
}

func CalculateBreakMinutes(shift model.Shifts) int {
	if shift.Break {
		// คำนวณเวลาพักโดยใช้ BreakOut และ BreakIn
		breakDuration := shift.BreakIn.Sub(shift.BreakOut)
		return int(breakDuration.Minutes())
	}
	return 0
}

func CalculateIsNightShift(shift model.Shifts) bool {
	// กรณีของคุณใช้ TIME อย่างเดียว (ไม่มีวันที่)
	// ถ้าเวลาเลิกงานน้อยกว่าเวลาเข้างาน แปลว่ากะนี้ข้ามวัน (เช่น 20:00 -> 04:00)
	return shift.EndTime.Before(shift.StartTime)
}

func CalculateLivingCost(shift model.Shifts) float64 {
	name := strings.TrimSpace(shift.ShiftName)

	// ถ้ามีคำว่า "ไม่มีค่าครองชีพ" ให้ถือว่า 0 ทันที
	if strings.Contains(name, "ไม่มีค่าครองชีพ") {
		return 0
	}

	// ต้องมีคำว่า "ค่าครองชีพ" ถึงจะพยายามดึงตัวเลข
	if !strings.Contains(name, "ค่าครองชีพ") {
		return 0
	}

	// สร้าง regex pattern สำหรับค้นหา "ตัวเลข" จากข้อความโดยตรงหลังคำว่า "ค่าครองชีพ"
	re := regexp.MustCompile(`ค่าครองชีพ\s*([0-9]+(?:\.[0-9]+)?)`)
	// หาตัวเลขใน name และเก็บไว้ใน match ถ้าเจอ
	match := re.FindStringSubmatch(name)
	if len(match) < 2 {
		return 0
	}

	val, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0
	}

	return val
}
