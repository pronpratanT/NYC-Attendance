package service

import (
	"log"
	"sync"

	"github.com/pronpratanT/leave-system/internal/attendance-service/model"
	"github.com/pronpratanT/leave-system/internal/attendance-service/repository"
)

type AttendanceService struct {
	CloudRepo *repository.CloudtimeRepository
	AppRepo   *repository.AttendanceRepository
}

func NewAttendanceService(cloudRepo *repository.CloudtimeRepository, appRepo *repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{
		CloudRepo: cloudRepo,
		AppRepo:   appRepo,
	}
}

func (s *AttendanceService) SyncFullLoad() error {

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
		s.syncRange(minBH, mid)
	}()

	go func() {
		defer wg.Done()
		// worker 2: [mid+1 .. maxBH] ป้องกันซ้ำกับ mid ของ worker แรก
		if mid+1 <= maxBH {
			s.syncRange(mid+1, maxBH)
		}
	}()

	wg.Wait()

	return nil
}

func (s *AttendanceService) syncRange(startBH, endBH int64) {

	batchSize := 3000
	// เริ่มจาก startBH-1 เพื่อให้เงื่อนไข bh > lastBH ครอบคลุม record แรกสุด (bh == startBH)
	lastBH := startBH - 1

	for {

		cloudRecords, err := s.CloudRepo.GetBatchByBHRange(lastBH, endBH, batchSize)
		if err != nil {
			log.Println("Fetch error:", err)
			return
		}

		if len(cloudRecords) == 0 {
			break
		}

		var insertData []model.Attendance

		for _, r := range cloudRecords {
			insertData = append(insertData, model.Attendance{
				BH:          r.BH,
				UserNo:      r.UserNo,
				SJ:          r.SJ,
				MC:          r.MC,
				UserLName:   r.UserLName,
				DepNo:       r.DepNo,
				UserDep:     r.UserDep,
				UserDepName: r.UserDepName,
				UserType:    r.UserType,
				UserCard:    r.UserCard,
			})
		}

		err = s.AppRepo.BulkInsert(insertData)
		if err != nil {
			log.Println("Insert error:", err)
			return
		}

		lastBH = cloudRecords[len(cloudRecords)-1].BH
	}
}
