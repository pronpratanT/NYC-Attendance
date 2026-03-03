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

func AttendanceRoutes(r *gin.RouterGroup, h *AttendanceHandler) {
	r.GET("/attendance-logs", h.GetAttendanceLogs)
	r.GET("/attendance-daily", h.GetAttendanceDaily)
	r.GET("/attendance-daily/by-employee/:employee_id", h.GetAttendanceDailyByEmployeeID)
	r.GET("/attendance-daily/by-employee/:employee_id/:start_date/:end_date", h.GetAttendanceDailyByEmployeeIDAndDate)
	r.GET("/attendance-daily/by-date/:start_date/:end_date", h.GetAttendanceDailyByDate)
}
