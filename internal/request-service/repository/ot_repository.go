package repository

import "gorm.io/gorm"

type OTRepository struct {
	DB *gorm.DB
}

func NewOTRepository(db *gorm.DB) *OTRepository {
	return &OTRepository{DB: db}
}
