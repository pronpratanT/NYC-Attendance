package service

import (
	"log"
	"sync"
	"time"

	"hr-program/internal/request-service/model"
)

func (s *RequestService) SyncFullLoad() error {
	minID, maxID, err := s.EconsRepo.GetMinMaxOTDocumentID()
	if err != nil {
		return err
	}

	mid := minID + (maxID-minID)/2

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// worker 1: [minID .. mid]
		s.syncRange(minID, mid)
	}()

	go func() {
		defer wg.Done()
		// worker 2: [mid+1 .. maxID] ป้องกันซ้ำกับ mid ของ worker แรก
		if mid+1 <= maxID {
			s.syncRange(mid+1, maxID)
		}
	}()

	wg.Wait()

	return nil
}

func (s *RequestService) syncRange(startID, endID int64) {

	batchSize := 3000
	// เริ่มจาก startID-1 เพื่อให้เงื่อนไข id > lastID ครอบคลุม record แรกสุด (id == startID)
	lastID := startID - 1

	for {
		records, err := s.EconsRepo.GetBatchOTByDocumentIDRange(lastID, endID, batchSize)
		if err != nil {
			log.Println("Fetch OT records error:", err)
			return
		}

		if len(records) == 0 {
			break
		}

		var insertData []model.OTlogs

		for _, r := range records {
			// Normalize date/time strings from ECONS to match Postgres column types
			dateStr := r.Date
			if t, err := time.Parse(time.RFC3339, r.Date); err == nil {
				dateStr = t.Format("2006-01-02")
			}

			startStr := r.StartOT
			if t, err := time.Parse(time.RFC3339, r.StartOT); err == nil {
				startStr = t.Format("15:04:05")
			}

			stopStr := r.StopOT
			if t, err := time.Parse(time.RFC3339, r.StopOT); err == nil {
				stopStr = t.Format("15:04:05")
			}

			insertData = append(insertData, model.OTlogs{
				ID:           r.ID,
				HRCheck:      r.HRCheck,
				Sequence:     r.Sequence,
				Department:   r.Department,
				Dep:          r.Dep,
				ShiftOT:      r.ShiftOT,
				TypeOT:       r.TypeOT,
				Date:         dateStr,
				AB:           r.AB,
				EmployeeCode: r.EmployeeCode,
				StartOT:      startStr,
				StopOT:       stopStr,
				WorkOT:       r.WorkOT,
				Approve:      r.Approve,
				RequestAP:    r.RequestAP,
				RequestTap:   r.RequestTap,
				ChiefAP:      r.ChiefAP,
				ChiefTap:     r.ChiefTap,
				ManagerAP:    r.ManagerAP,
				ManagerTap:   r.ManagerTap,
				HRAP:         r.HRAP,
				HRTap:        r.HRTap,
				DeleteName:   r.DeleteName,
				Deletetime:   r.Deletetime,
				CreateDate:   r.CreateDate,
			})
		}

		err = s.AppRepo.BulkInsert(insertData)
		if err != nil {
			log.Println("Batch upsert OT error:", err)
			return
		}

		lastID = records[len(records)-1].ID
	}
}
