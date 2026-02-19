package repository

import (
	"hr-program/internal/attendance-service/model"

	"gorm.io/gorm"
)

type CloudtimeRepository struct {
	DB *gorm.DB
}

func NewCloudtimeRepository(db *gorm.DB) *CloudtimeRepository {
	return &CloudtimeRepository{DB: db}
}

// หา min และ max BH
func (r *CloudtimeRepository) GetMinMaxBH() (int64, int64, error) {

	var minBH, maxBH int64

	row := r.DB.Model(&model.CloudtimeAttendance{}).
		Select("MIN(bh), MAX(bh)").
		Row()
	err := row.Scan(&minBH, &maxBH)

	return minBH, maxBH, err
}

// ดึง batch ตามช่วงเวลา + keyset
func (r *CloudtimeRepository) GetBatchByBHRange(
	lastBH int64,
	endBH int64,
	limit int,
) ([]model.CloudtimeAttendance, error) {

	var records []model.CloudtimeAttendance

	err := r.DB.
		Where("bh > ? AND bh <= ?", lastBH, endBH).
		Order("bh ASC").
		Limit(limit).
		Find(&records).Error

	return records, err
}
