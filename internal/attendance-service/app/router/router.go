package router

import (
	"hr-program/internal/attendance-service/handler"

	"github.com/gin-gonic/gin"
)

func AttendanceRouter(r *gin.Engine, h *handler.AttendanceHandler) *gin.Engine {
	api := r.Group("/api")
	handler.AttendanceRoutes(api.Group("/attendance"), h)

	return r
}
