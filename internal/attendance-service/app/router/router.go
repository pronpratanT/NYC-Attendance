package router

import (
	"hr-program/internal/attendance-service/handler"

	"github.com/gin-gonic/gin"
)

func AttendanceRouter(h *handler.AttendanceHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/attendance")
	{
		api.GET("/attendance-logs", h.GetAttendanceLogs)
		api.GET("/attendance-daily", h.GetAttendanceDaily)
		// แยก path ชัดเจนเพื่อไม่ให้ Gin conflict ระหว่าง wildcard ตัวแรก
		api.GET("/attendance-daily/by-employee/:employee_id", h.GetAttendanceDailyByEmployeeID)
		api.GET("/attendance-daily/by-employee/:employee_id/:start_date/:end_date", h.GetAttendanceDailyByEmployeeIDAndDate)
		api.GET("/attendance-daily/by-date/:start_date/:end_date", h.GetAttendanceDailyByDate)
	}

	return r
}
