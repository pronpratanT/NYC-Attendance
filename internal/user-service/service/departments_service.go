package service

import (
	model "hr-program/internal/user-service/model/departments"
	repository "hr-program/internal/user-service/repository/departments"
	"log"
	"sync"
)

type DepartmentsService struct {
	CloudtimeRepo *repository.CloudtimeDepartmentsRepository
	AppRepo       *repository.DepartmentsRepository
}

func NewDepartmentsService(cloudRepo *repository.CloudtimeDepartmentsRepository, appRepo *repository.DepartmentsRepository) *DepartmentsService {
	return &DepartmentsService{
		CloudtimeRepo: cloudRepo,
		AppRepo:       appRepo,
	}
}

func (s *DepartmentsService) SyncFullLoad() error {

	minDepNo, maxDepNo, err := s.CloudtimeRepo.GetMinMaxDepSerial()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// Sync full load ในช่วง [minDepNo .. maxDepNo]
		s.syncRange(minDepNo, maxDepNo)
	}()

	wg.Wait()

	return nil
}

func (s *DepartmentsService) syncRange(startDepNo, endDepNo int64) {

	batchSize := 3000
	// เริ่มจาก startDepNo-1 เพื่อให้เงื่อนไข dep_no > lastDepNo ครอบคลุม record แรกสุด (dep_no == startDepNo)
	lastDepNo := startDepNo - 1

	for {
		cloudRecords, err := s.CloudtimeRepo.GetBatchByDepSerialRange(lastDepNo, endDepNo, batchSize)
		if err != nil {
			log.Println("Fetch departments error:", err)
			return
		}

		if len(cloudRecords) == 0 {
			break
		}

		var insertData []model.Departments

		for _, r := range cloudRecords {
			// แปลงข้อมูลจาก CloudtimeDepartments เป็น Departments
			insertData = append(insertData, model.Departments{
				DepNo: r.DepNo,
				Name:  r.DepName,
			})
		}

		err = s.AppRepo.BulkInsert(insertData)
		if err != nil {
			log.Println("Insert departments error:", err)
			return
		}

		// อัปเดต lastDepNo เป็น dep_no ของ record สุดท้ายที่ดึงมา
		lastDepNo = cloudRecords[len(cloudRecords)-1].DepSerial
	}
}
