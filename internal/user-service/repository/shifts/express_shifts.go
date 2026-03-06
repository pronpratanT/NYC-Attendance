package repository

import (
	model "hr-program/shared/models/users"

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

func (r *SqlExpressShiftRepository) GetUserBplus() ([]model.SQLExpressUser, error) {
	var usr []model.SQLExpressUser
	err := r.DB.Find(&usr).Error
	return usr, err
}

func (r *SqlExpressShiftRepository) GetMasterKey() ([]model.SQLExpressMasterKey, error) {
	var mk []model.SQLExpressMasterKey
	err := r.DB.Where("TMR_DATE = ?", "2026-03-04T00:00:00Z").Find(&mk).Error
	return mk, err
}

// GetLatestUserShifts เดิมใช้ดึง N แถวล่าสุด ตอนนี้ปรับให้ดึงทั้งหมด โดยเรียงตาม SF_KEY จากมากไปน้อย
// func (r *SqlExpressShiftRepository) GetLatestUserShifts(limit int) ([]model.SQLExpressUser, error) {
// 	var usr []model.SQLExpressUser

// 	query := fmt.Sprintf("SELECT * FROM EMPFILE")
// 	err := r.DB.Raw(query).Scan(&usr).Error
// 	return usr, err
// }

// func (r *SqlExpressShiftRepository) GetLatestMaster() ([]model.SQLExpressMasterKey, error) {
// 	var mk []model.SQLExpressMasterKey

// 	query := fmt.Sprintf("SELECT * FROM TMRESULT WHERE TMR_DATE = '2026-03-04T00:00:00Z'")
// 	err := r.DB.Raw(query).Scan(&mk).Error
// 	return mk, err
// }

// func (r *SqlExpressShiftRepository) GetLatestShifts(limit int) ([]model.SQLExpressShifts, error) {
// 	var shifts []model.SQLExpressShifts

// 	query := fmt.Sprintf("SELECT * FROM TMSHIFT")
// 	err := r.DB.Raw(query).Scan(&shifts).Error
// 	return shifts, err
// }
