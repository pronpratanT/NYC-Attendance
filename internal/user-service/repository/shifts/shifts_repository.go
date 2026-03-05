package repository

import (
	model "hr-program/shared/models/users"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShiftsRepository struct {
	DB *gorm.DB
}

func NewShiftsRepository(db *gorm.DB) *ShiftsRepository {
	return &ShiftsRepository{DB: db}
}

func (r *ShiftsRepository) BulkInsertShifts(data []model.Shifts) error {
	return r.DB.
		Clauses(clause.OnConflict{
			// ไม่ระบุ column ให้ Postgres สร้าง "ON CONFLICT DO NOTHING" ทั่วไป
			// ให้จัดการตาม primary key / constraints ที่มีอยู่
			DoNothing: true,
		}).
		CreateInBatches(data, len(data)).Error
}
