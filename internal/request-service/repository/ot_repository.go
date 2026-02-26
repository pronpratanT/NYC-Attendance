package repository

import (
	model "hr-program/internal/request-service/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OTRepository struct {
	DB *gorm.DB
}

func NewOTRepository(db *gorm.DB) *OTRepository {
	return &OTRepository{DB: db}
}

func (r *OTRepository) BulkInsert(data []model.OTlogs) error {
	// Use a safe batch size to avoid Postgres 65535-parameter limit
	const batchSize = 500

	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).
		CreateInBatches(data, batchSize).Error
}

func (r *OTRepository) GetOTlogs() ([]model.OTlogs, error) {
	sqlDB, err := r.DB.DB()
	if err != nil {
		log.Println("Failed to get raw DB connection:", err)
		return nil, err
	}
	rows, err := sqlDB.Query("SELECT id, hr_check, sequence, department, dep, shift_ot, type_ot, date, ab, employee_code, start_ot, stop_ot, work_ot, approve, request_ap, request_tap, chief_ap, chief_tap, manager_ap, manager_tap, hr_ap, hr_tap, delete_name, delete_time, create_date FROM ot_logs")
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var otlogs []model.OTlogs
	for rows.Next() {
		var log model.OTlogs
		err := rows.Scan(&log.ID, &log.HRCheck, &log.Sequence, &log.Department, &log.Dep, &log.ShiftOT, &log.TypeOT, &log.Date, &log.AB, &log.EmployeeCode, &log.StartOT, &log.StopOT, &log.WorkOT, &log.Approve, &log.RequestAP, &log.RequestTap, &log.ChiefAP, &log.ChiefTap, &log.ManagerAP, &log.ManagerTap, &log.HRAP, &log.HRTap, &log.DeleteName, &log.Deletetime, &log.CreateDate)
		if err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}
		otlogs = append(otlogs, log)
	}
	return otlogs, nil
}

func (r *OTRepository) SaveOTDoc(docs []model.OTDoc) error {
	// ใช้ batch size เล็กลงเพื่อเลี่ยง limit 65535 parameters ของ Postgres
	const batchSize = 500
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sequence"}},
			DoNothing: true,
		}).
		CreateInBatches(docs, batchSize).Error
}

func (r *OTRepository) SaveOTDetails(details []model.OTDetail) error {
	// ใช้ batch size เล็กลงเพื่อเลี่ยง limit 65535 parameters ของ Postgres
	const batchSize = 500
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "source_log_id"}},
			DoNothing: true,
		}).
		CreateInBatches(details, batchSize).Error
}

func (r *OTRepository) GetOTDocsBySequences(sequences []int64) ([]model.OTDoc, error) {
	var docs []model.OTDoc
	if len(sequences) == 0 {
		return docs, nil
	}
	if err := r.DB.Where("sequence IN ?", sequences).Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}
