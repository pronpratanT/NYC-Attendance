package repository

import (
	model "hr-program/internal/request-service/model"

	"gorm.io/gorm"
)

type EconsRepository struct {
	// DB connection หรือ client สำหรับเชื่อมต่อกับ ECONS
	DB *gorm.DB
}

func NewEconsRepository(db *gorm.DB) *EconsRepository {
	return &EconsRepository{DB: db}
}

func (r *EconsRepository) GetMinMaxOTDocumentID() (int64, int64, error) {
	var minID, maxID int64

	row := r.DB.Model(&model.OTEcons{}).
		Select("MIN(id), MAX(id)").
		Row()
	err := row.Scan(&minID, &maxID)

	return minID, maxID, err
}

func (r *EconsRepository) GetBatchOTByDocumentIDRange(
	lastID int64,
	endID int64,
	limit int,
) ([]model.OTEcons, error) {
	var records []model.OTEcons

	err := r.DB.
		Where("id > ? AND id <= ?", lastID, endID).
		Order("id ASC").
		Limit(limit).
		Find(&records).Error

	return records, err
}

func (r *EconsRepository) GetAllHolidays() ([]model.HolidayEcons, error) {
	var records []model.HolidayEcons

	err := r.DB.
		Order("date ASC").
		Find(&records).Error

	return records, err
}
