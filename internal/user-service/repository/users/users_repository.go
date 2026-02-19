package repository

import (
	model "hr-program/internal/user-service/model/users"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) BulkInsert(data []model.Users) error {
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "employee_id"}},
			DoNothing: true,
		}).
		CreateInBatches(data, len(data)).Error
}
