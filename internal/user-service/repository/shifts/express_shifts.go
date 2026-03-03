package repository

import (
	// model "hr-program/shared/models/users"

	"log"

	"gorm.io/gorm"
)

type SqlExpressShiftRepository struct {
	DB *gorm.DB
}

func NewSQLExpressShiftRepository(db *gorm.DB) *SqlExpressShiftRepository {
	return &SqlExpressShiftRepository{DB: db}
}

// func (r *SqlExpressShiftRepository) GetAllShifts() ([]model.SQLExpressShifts, error) {
// 	var shifts []model.SQLExpressShifts
// 	err := r.DB.Find(&shifts).Error
// 	return shifts, err
// }

func (r *SqlExpressShiftRepository) DebugDescribeTMSHIFT() error {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return err
	}

	rows, err := sqlDB.Query("SELECT TOP 1 * FROM PERSONALINFO")
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for _, c := range cols {
		log.Printf("Column: %s, DB Type: %s", c.Name(), c.DatabaseTypeName())
	}
	return nil
}
