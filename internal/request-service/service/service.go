package service

import "hr-program/internal/request-service/repository"

type RequestService struct {
	AppRepo   *repository.OTRepository
	EconsRepo *repository.EconsRepository
	UserRepo  repository.UserRepositoryInterface
}

func NewRequestService(appRepo *repository.OTRepository, econsRepo *repository.EconsRepository, userRepo repository.UserRepositoryInterface) *RequestService {
	return &RequestService{
		AppRepo:   appRepo,
		EconsRepo: econsRepo,
		UserRepo:  userRepo,
	}
}
