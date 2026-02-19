package repository

import (
	model "hr-program/internal/user-service/model/users"

	"gorm.io/gorm"
)

type CloudtimeUserRepository struct {
	DB *gorm.DB
}

func NewCloudtimeUserRepository(db *gorm.DB) *CloudtimeUserRepository {
	return &CloudtimeUserRepository{DB: db}
}

func (r *CloudtimeUserRepository) GetMinMaxUserSerial() (int64, int64, error) {
	var minID, maxID int64

	row := r.DB.Model(&model.CloudtimeUser{}).
		Select("MIN(user_serial), MAX(user_serial)").
		Row()
	err := row.Scan(&minID, &maxID)

	return minID, maxID, err
}

func (r *CloudtimeUserRepository) GetBatchByUserSerialRange(
	lastID int64,
	endID int64,
	limit int,
) ([]model.CloudtimeUser, error) {
	var records []model.CloudtimeUser

	err := r.DB.
		Where("user_serial > ? AND user_serial <= ?", lastID, endID).
		Order("user_serial ASC").
		Limit(limit).
		Find(&records).Error

	return records, err
}
