package handler

import (
	"hr-program/internal/attendance-service/service"
)

type AttendanceHandler struct {
	Service *service.AttendanceService
}

func NewAttendanceHandler(s *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{Service: s}
}
