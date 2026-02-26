package service

import (
	"time"

	"hr-program/internal/request-service/model"
)

func (s *RequestService) GenerateAndSaveOT() error {
	docs, details, err := s.OTLogsProcessing()
	if err != nil {
		return err
	}
	if len(docs) == 0 || len(details) == 0 {
		return nil
	}

	// 1) เซฟหัวเอกสาร
	if err := s.AppRepo.SaveOTDoc(docs); err != nil {
		return err
	}

	// 2) ดึง OTDoc จาก DB อีกรอบ เพื่อให้ได้ id ที่ถูกต้อง ทั้งของ record ใหม่และ record เดิม
	sequences := make([]int64, 0, len(docs))
	for _, d := range docs {
		sequences = append(sequences, d.Sequence)
	}
	savedDocs, err := s.AppRepo.GetOTDocsBySequences(sequences)
	if err != nil {
		return err
	}

	// 3) map sequence -> id ที่อยู่ใน DB จริง
	seqToID := make(map[int64]uint, len(savedDocs))
	for _, d := range savedDocs {
		seqToID[d.Sequence] = d.ID
	}

	// 4) ใส่ OTDocID ให้ detail แต่ละแถว ถ้าไม่มี doc สำหรับ sequence นั้นจะข้ามไป
	filteredDetails := make([]model.OTDetail, 0, len(details))
	for i := range details {
		id, ok := seqToID[details[i].Sequence]
		if !ok {
			continue
		}
		details[i].OTDocID = id
		filteredDetails = append(filteredDetails, details[i])
	}
	if len(filteredDetails) == 0 {
		return nil
	}

	// 5) เซฟรายละเอียด
	return s.AppRepo.SaveOTDetails(filteredDetails)
}

func (s *RequestService) OTLogsProcessing() ([]model.OTDoc, []model.OTDetail, error) {
	otLogs, err := s.AppRepo.GetOTlogs()
	if err != nil {
		return nil, nil, err
	}

	// เตรียม map แปลง employee_code -> user_id จาก user-service
	employeeSet := make(map[string]struct{})
	for _, ot := range otLogs {
		if ot.EmployeeCode == "" {
			continue
		}
		employeeSet[ot.EmployeeCode] = struct{}{}
	}

	employeeCodes := make([]string, 0, len(employeeSet))
	for code := range employeeSet {
		employeeCodes = append(employeeCodes, code)
	}

	userMap, err := s.UserRepo.GetUserIDMapByEmployeeIDs(employeeCodes)
	if err != nil {
		return nil, nil, err
	}

	// group by sequence to build OTDoc, then flatten employees into OTDetail
	docsMap := make(map[int64]*model.OTDoc)
	details := make([]model.OTDetail, 0, len(otLogs))

	for _, r := range otLogs {
		// สนใจเฉพาะรายการที่ HR_CHECK = "APPROVE" เท่านั้น
		if r.HRCheck != "APPROVE" {
			continue
		}

		// parse date and times from strings in ot_logs
		var parsedDate time.Time
		// รองรับทั้งรูปแบบ "2006-01-02" และ RFC3339 เช่น "2006-01-02T00:00:00Z"
		parsedDate, err = time.Parse("2006-01-02", r.Date)
		if err != nil {
			// ลอง parse แบบ RFC3339 แล้วตัดให้เหลือแค่วันที่
			if t, err2 := time.Parse(time.RFC3339, r.Date); err2 == nil {
				parsedDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			} else {
				return nil, nil, err
			}
		}
		startTime, err := time.Parse("15:04:05", r.StartOT)
		if err != nil {
			return nil, nil, err
		}
		stopTime, err := time.Parse("15:04:05", r.StopOT)
		if err != nil {
			return nil, nil, err
		}

		var userID int64
		if id, ok := userMap[r.EmployeeCode]; ok {
			userID = id
		}

		doc, ok := docsMap[r.Sequence]
		if !ok {
			chiefTap := r.ChiefTap
			managerTap := r.ManagerTap
			hrTap := r.HRTap

			doc = &model.OTDoc{
				Sequence:   r.Sequence,
				Date:       parsedDate,
				ShiftOT:    r.ShiftOT,
				Department: r.Department,
				Dep:        r.Dep,
				RequestAP:  r.RequestAP,
				RequestTap: r.RequestTap,
				HRCheck:    r.HRCheck,
				Approve:    r.Approve,
				ChiefAP:    r.ChiefAP,
				ChiefTap:   &chiefTap,
				ManagerAP:  r.ManagerAP,
				ManagerTap: &managerTap,
				HRAP:       r.HRAP,
				HRTap:      &hrTap,
			}
			docsMap[r.Sequence] = doc
		} else {
			// ensure doc-level date is consistent; if not, keep the first one
			_ = doc
		}

		detail := model.OTDetail{
			OTDocID:     0, // จะถูกเติมภายหลัง หลังจากบันทึก OTDoc แล้ว map จาก Sequence -> ID
			EmployeeID:  userID,
			TypeOT:      r.TypeOT,
			StartOT:     startTime,
			StopOT:      stopTime,
			WorkOT:      r.WorkOT,
			SourceLogID: r.ID,
			Sequence:    r.Sequence,
		}
		details = append(details, detail)
	}

	otDocs := make([]model.OTDoc, 0, len(docsMap))
	for _, d := range docsMap {
		otDocs = append(otDocs, *d)
	}

	return otDocs, details, nil
}
