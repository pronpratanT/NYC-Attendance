package repository

import (
	"fmt"
	model "hr-program/shared/models/users"
	"log"

	"gorm.io/gorm"
)

type SqlExpressShiftRepository struct {
	DB *gorm.DB
}

func NewSQLExpressShiftRepository(db *gorm.DB) *SqlExpressShiftRepository {
	return &SqlExpressShiftRepository{DB: db}
}

// GetAllShifts ดึง shift ทั้งหมดจาก TMSHIFT
func (r *SqlExpressShiftRepository) GetAllShifts() ([]model.SQLExpressShifts, error) {
	var shifts []model.SQLExpressShifts
	err := r.DB.Find(&shifts).Error
	return shifts, err
}

// GetLatestUserShifts เดิมใช้ดึง N แถวล่าสุด ตอนนี้ปรับให้ดึงทั้งหมด โดยเรียงตาม SF_KEY จากมากไปน้อย
func (r *SqlExpressShiftRepository) GetLatestUserShifts(limit int) ([]model.SQLExpressUser, error) {
	var usr []model.SQLExpressUser

	query := fmt.Sprintf("SELECT * FROM PERSONALINFO")
	err := r.DB.Raw(query).Scan(&usr).Error
	return usr, err
}

func (r *SqlExpressShiftRepository) GetUserRaw() ([]map[string]interface{}, error) {
	var rows []map[string]interface{}

	// ดึงไม่เกิน 10 แถวจาก TMRESULT (SQL Server ใช้ TOP)
	if err := r.DB.Raw("SELECT TOP 10 * FROM TMRESULT").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

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
