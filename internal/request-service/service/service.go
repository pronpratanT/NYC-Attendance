package service

import "hr-program/internal/request-service/repository"

type RequestService struct {
	OTRepo    *repository.OTRepository
	EconsRepo *repository.EconsRepository
}

func NewRequestService(otRepo *repository.OTRepository, econsRepo *repository.EconsRepository) *RequestService {
	return &RequestService{
		OTRepo:    otRepo,
		EconsRepo: econsRepo,
	}
}
