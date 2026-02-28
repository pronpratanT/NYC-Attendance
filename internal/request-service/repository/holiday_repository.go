package repository

import (
	model "hr-program/internal/request-service/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HolidayRepository struct {
	DB *gorm.DB
}

func NewHolidayRepository(db *gorm.DB) *HolidayRepository {
	return &HolidayRepository{DB: db}
}

func (r *HolidayRepository) BulkInsertHolidays(data []model.Holiday) error {
	// Use a safe batch size to avoid Postgres 65535-parameter limit
	const batchSize = 500

	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).
		CreateInBatches(data, batchSize).Error
}
