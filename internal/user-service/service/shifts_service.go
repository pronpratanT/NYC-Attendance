package service

import (
	model "hr-program/shared/models/users"
	"regexp"
	"strconv"
	"strings"
)

func (s *UserService) GetLatestShifts(limit int) ([]model.SQLExpressUser, error) {
	return s.SQLExpressRepo.GetLatestUserShifts(limit)
}

// GetUserRaw ดึงข้อมูลทุก column จาก PERSONALINFO แบบไม่ผูกกับ struct
// คืนค่าเป็น slice ของ map[columnName]value
func (s *UserService) GetUserRaw() ([]map[string]interface{}, error) {
	return s.SQLExpressRepo.GetUserRaw()
}

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
