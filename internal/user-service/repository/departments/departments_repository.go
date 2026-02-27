package repository

import (
	model "hr-program/internal/user-service/model/departments"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DepartmentsRepository struct {
	DB *gorm.DB
}

func NewDepartmentsRepository(db *gorm.DB) *DepartmentsRepository {
	return &DepartmentsRepository{DB: db}
}

func (r *DepartmentsRepository) BulkInsert(data []model.Departments) error {
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "dep_no"}},
			DoNothing: true,
		}).
		CreateInBatches(data, len(data)).Error
}

// GetDepartmentsIDMap คืนค่า map จาก dep_no (string) ไปเป็น department id (int64)
// ใช้สำหรับ map รหัสแผนกจาก Cloudtime เข้ากับตาราง departments ในแอป
func (r *DepartmentsRepository) GetDepartmentsIDMap(depNos []string) (map[string]int64, error) {
	var departments []model.Departments

	err := r.DB.
		Select("id", "dep_no").
		Where("dep_no IN ?", depNos).
		Find(&departments).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int64, len(departments))
	for _, d := range departments {
		result[d.DepNo] = d.ID
	}

	return result, nil
}
