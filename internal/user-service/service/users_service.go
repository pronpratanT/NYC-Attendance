package service

import (
	model "hr-program/shared/models/users"
	"log"
	"strings"
	"sync"
)

func (s *UserService) SyncFullLoadUsers() error {

	minID, maxID, err := s.CloudtimeUserRepo.GetMinMaxUserSerial()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// Sync full load ในช่วง [minID .. maxID]
		s.syncRangeUsers(minID, maxID)
	}()

	wg.Wait()

	return nil
}

func (s *UserService) syncRangeUsers(startID, endID int64) {

	batchSize := 3000
	// เริ่มจาก startID-1 เพื่อให้เงื่อนไข id > lastID ครอบคลุม record แรกสุด (id == startID)
	lastID := startID - 1

	for {
		cloudRecords, err := s.CloudtimeUserRepo.GetBatchByUserSerialRange(lastID, endID, batchSize)
		if err != nil {
			log.Println("Fetch users error:", err)
			return
		}

		if len(cloudRecords) == 0 {
			break
		}

		var insertData []model.Users

		for _, r := range cloudRecords {
			// แปลงชื่อเต็มเป็นชื่อและนามสกุล
			fn, ln := splitFullName(r.UserLname)

			insertData = append(insertData, model.Users{
				EmployeeID:   r.UserNo,
				Password:     r.UserNo,
				DepartmentID: r.UserDep,
				FName:        fn,
				LName:        ln,
				IsActive:     true,
				Workday:      r.UserWorkday,
				BirthDate:    r.UserBirthday,
			})
		}

		err = s.AppRepo.BulkInsert(insertData)
		if err != nil {
			log.Println("Insert user error:", err)
			return
		}

		lastID = cloudRecords[len(cloudRecords)-1].UserSerial
	}
}

func splitFullName(full string) (first, last string) {
	// TrimSpcae ในการตัดช่องว่างส่วนเกิน หน้าและหลัง ข้อความ
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}

	// Field จะตัดข้อความโดยใช้ช่องว่างเป็นตัวแบ่ง และคืนค่าเป็น slice ของ string
	parts := strings.Fields(full)
	// ถ้ามีคำเดียว ให้ถือว่าเป็นชื่อแรก และชื่อสุดท้ายเป็นค่าว่าง
	if len(parts) == 1 {
		return parts[0], ""
	}

	// ถ้ามีครบทั้งชื่อและสกุล parts[0] = first parts[1:] = last [1:] คือคำที่เหลือทั้งหมด
	return parts[0], strings.Join(parts[1:], " ")
}
