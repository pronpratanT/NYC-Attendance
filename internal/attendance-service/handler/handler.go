package handler

import (
	"hr-program/internal/attendance-service/service"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Service *service.AttendanceService
}

func NewAttendanceHandler(s *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{Service: s}
}

// ฟังก์ชันช่วย register routes
func RegisterAttendanceRoutes(r *gin.Engine, h *AttendanceHandler) {
	r.GET("/attendance/logs", h.GetAttendanceLogs)
}
