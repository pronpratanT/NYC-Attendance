package repository

import (
	model "hr-program/internal/user-service/model/departments"

	"gorm.io/gorm"
)

type CloudtimeDepartmentsRepository struct {
	DB *gorm.DB
}

func NewCloudtimeDepartmentsRepository(db *gorm.DB) *CloudtimeDepartmentsRepository {
	return &CloudtimeDepartmentsRepository{DB: db}
}

func (r *CloudtimeDepartmentsRepository) GetMinMaxDepSerial() (int64, int64, error) {
	var minID, maxID int64

	row := r.DB.Model(&model.CloudtimeDepartments{}).
		Select("MIN(dep_serial), MAX(dep_serial)").
		Row()
	err := row.Scan(&minID, &maxID)

	return minID, maxID, err
}

func (r *CloudtimeDepartmentsRepository) GetBatchByDepSerialRange(
	lastID int64,
	endID int64,
	limit int,
) ([]model.CloudtimeDepartments, error) {
	var records []model.CloudtimeDepartments

	err := r.DB.
		Where("dep_serial > ? AND dep_serial <= ?", lastID, endID).
		Order("dep_serial ASC").
		Limit(limit).
		Find(&records).Error

	return records, err
}
