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
	}

	return r
}
