package repository

import "gorm.io/gorm"

type ShiftsRepository struct {
	DB *gorm.DB
}

func NewShiftsRepository(db *gorm.DB) *ShiftsRepository {
	return &ShiftsRepository{DB: db}
}
