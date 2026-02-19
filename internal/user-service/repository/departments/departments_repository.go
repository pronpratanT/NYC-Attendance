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
