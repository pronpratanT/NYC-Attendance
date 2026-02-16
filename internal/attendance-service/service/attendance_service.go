package service

import (
	"context"
	"hr-backend/internal/attendance-service/repository"

	"time"

	"github.com/pronpratanT/leave-system/internal/attendance-service/model"
	"github.com/pronpratanT/leave-system/internal/attendance-service/repository"
)

type AttendanceService struct {
	Repo *repository.AttendanceRepository
}

func NewAttendanceService(r *repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{Repo: r}
}

func (s *AttendanceService) SyncTest(ctx context.Context) error {
	now := time.Now()

	data := model.Attendance{
		BH:          1,
		UserSerial:  1001,
		UserNo:      "EMP001",
		UserLName:   "John Doe",
		DepNo:       "D001",
		UserDep:     10,
		UserDepName: "HR",
		UserType:    1,
		UserCard:    "CARD001",
		SJ:          now,
		Iden:        "IDENT001",
		FX:          0,
	}

	return s.Repo.Insert(ctx, data)
}
