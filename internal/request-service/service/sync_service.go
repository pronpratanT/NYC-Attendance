package service

import (
	"log"
	"sync"

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

		var insertData []model.OT

		for _, r := range records {
			insertData = append(insertData, model.OT{
				ID:              r.ID,
				OTDocumentID:    r.Sequence,
				TYPE_OT:         r.TypeOT,
				OTRequesterName: r.RequestAP,
				ChieftainName:   r.ChiefAP,
				ManagerName:     r.ManagerAP,
				HRApproveName:   r.HRAP,
				HRCheckStatus:   r.HRCheck,
				ApproveStatus:   r.Approve,
				CreatedAt:       r.CreateDate,
			}),
			insertData = append(insertData, model.OTDetail{
				ID:           r.ID,
				OTDocumentID: r.Sequence,
				TYPE_OT:      r.TypeOT,
				DepartmentID: r.Dep, // ต้องแก้ไข ยต้องนำไป map id ก่อน
				ShiftID:      r.ShiftOT, // ต้องแก้ไข ยต้องนำไป map id ก่อน
				EmployeeID:   r.EmployeeCode, // ต้องแก้ไข ยต้องนำไป map id ก่อน
				OTDate:       r.Date,
				OTStart:      r.StartOT,
				OTEnd:        r.StopOT,
				CreatedAt:    r.CreateDate,
			})
		}

		err := s.AppRepo.BulkInsert(insertData)
		if err != nil {
			log.Println("Batch upsert OT error:", err)
			return
		}

		lastID = records[len(records)-1].ID
	}
}
