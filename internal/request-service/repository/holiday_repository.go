package repository

import (
	model "hr-program/shared/models/request"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HolidayRepository struct {
	DB *gorm.DB
}

func NewHolidayRepository(db *gorm.DB) *HolidayRepository {
	return &HolidayRepository{DB: db}
}

type HolidayRepositoryInterface interface {
	GetHolidays() ([]model.Holiday, error)
	GetHolidayByDate(date string) ([]model.Holiday, error)
	GetHolidayByDateRange(startDate, endDate string) ([]model.Holiday, error)
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

func (r *HolidayRepository) GetHolidays() ([]model.Holiday, error) {
	var holidays []model.Holiday
	err := r.DB.Find(&holidays).Error
	return holidays, err
}

func (r *HolidayRepository) GetHolidayByDate(date string) ([]model.Holiday, error) {
	var holiday []model.Holiday
	err := r.DB.Where("date = ?", date).Find(&holiday).Error
	return holiday, err
}

func (r *HolidayRepository) GetHolidayByDateRange(startDate, endDate string) ([]model.Holiday, error) {
	var holidays []model.Holiday
	err := r.DB.Where("date <= ? AND date >= ?", endDate, startDate).Find(&holidays).Error
	return holidays, err
}
