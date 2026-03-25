package handler

import (
	"hr-program/internal/attendance-service/service"
	"hr-program/shared/middleware"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Service *service.AttendanceService
}

func NewAttendanceHandler(s *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{Service: s}
}

func AttendanceRoutes(r *gin.RouterGroup, h *AttendanceHandler) {
	// public := r.Group("")
	// public.GET("/attendance-logs/date-range", h.GetAttendanceLogsByDateRange)

	protected := r.Group("")
	protected.Use(middleware.JWTAuth())
	protected.GET("/attendance-logs", h.GetAttendanceLogs)
	protected.GET("/attendance-daily", h.GetAttendanceDaily)
	protected.GET("/attendance-daily/by-employee/:employee_id", h.GetAttendanceDailyByEmployeeID)
	protected.GET("/attendance-daily/by-employee/:employee_id/:start_date/:end_date", h.GetAttendanceDailyByEmployeeIDAndDateRange)
	protected.GET("/attendance-daily/by-date/:start_date/:end_date", h.GetAttendanceDailyByDateRange)
	protected.GET("/attendance-logs/date-range", h.GetAttendanceLogsByDateRange)
	protected.GET("/attendance-logs/export/txt", h.ExportAttendanceLogsTXTByDateRange)
}
